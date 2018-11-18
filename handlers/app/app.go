package app

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

//Root just make sure we can hits the redis
func Root(c *gin.Context) {
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}
	pong, err := redisConn.Ping().Result()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Cannot ping Redis:",
			"error":   err,
		})
	}
	fmt.Println(pong, err)
	c.JSON(200, gin.H{
		"message": "This is root",
	})
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

	pipe := redisConn.Pipeline()
	score := pipe.ZScore(c.Param("set"), c.Param("member"))
	rank := pipe.ZRank(c.Param("set"), c.Param("member"))
	_, err := pipe.Exec()

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Pipe failed:",
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "OK",
			"board":   c.Param("set"),
			"member":  c.Param("member"),
			"score":   score.Val(),
			"rank":    rank.Val(),
		})
	}
}
