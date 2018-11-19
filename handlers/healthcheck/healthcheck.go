package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func HealthCheck(c *gin.Context) {
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}
	_, err := redisConn.Ping().Result()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Cannot ping Redis:",
			"error":   err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	}
}
