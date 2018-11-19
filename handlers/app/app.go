package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

//MemberRank struct to get member information
type MemberRank struct {
	score float64
	rank  int64
}

//Root just make sure we can hits the redis
func Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "This is root - Please see the docs",
	})
}

func getrank(redisConn *redis.Client, set string, member string) (rank MemberRank, rerr error) {

	pipe := redisConn.Pipeline()
	mscore := pipe.ZScore(set, member)
	mrank := pipe.ZRevRank(set, member)
	_, rerr = pipe.Exec()
	rank = MemberRank{rank: mrank.Val(), score: mscore.Val()}
	return
}

//Incr increments the leaderboard and returns current score
func Incr(c *gin.Context) {
	var incby float64
	s, p := strconv.ParseFloat(c.Param("count"), 64)
	if p == nil {
		incby = s
	} else {
		incby = 1
	}
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}

	incErr := redisConn.ZIncrBy(c.Param("set"), incby, c.Param("member")).Err()
	if incErr != nil {
		c.JSON(500, gin.H{
			"message": "Unable to set",
			"error":   incErr,
		})
	}

	m, err := getrank(redisConn, c.Param("set"), c.Param("member"))

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Rank fetch failed:",
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"board":  c.Param("set"),
			"member": c.Param("member"),
			"score":  m.score,
			"rank":   m.rank,
		})
	}
}

//GetRank for a member of a set
func GetRank(c *gin.Context) {
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}

	m, err := getrank(redisConn, c.Param("set"), c.Param("member"))

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Rank fetch failed:",
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"board":  c.Param("set"),
			"member": c.Param("member"),
			"score":  m.score,
			"rank":   m.rank,
		})
	}
}

//ShowBoard returns the leaderboard as a JSON object
func ShowBoard(c *gin.Context) {
	var entryCount int64
	s, p := strconv.ParseInt(c.Param("count"), 10, 64)
	if p == nil {
		entryCount = s - 1
	} else {
		entryCount = -1
	}
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}

	board := redisConn.ZRevRangeWithScores(c.Param("set"), 0, entryCount)

	if board.Err() != nil {
		c.JSON(500, gin.H{
			"message": "Board fetch failed:",
			"error":   board.Err(),
		})
	} else {
		c.JSON(200, gin.H{
			"board":   c.Param("set"),
			"leaders": board.Val(),
		})
	}
}
