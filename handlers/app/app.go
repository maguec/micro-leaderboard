package app

import (
	"fmt"

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
