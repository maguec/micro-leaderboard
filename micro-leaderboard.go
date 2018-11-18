package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/maguec/micro-leaderboard/handlers/app"
	"github.com/maguec/micro-leaderboard/handlers/healthcheck"
	"github.com/shokunin/contrib/ginrus"
	"github.com/sirupsen/logrus"
)

// APIMiddleware will add the redis connection to the context
func APIMiddleware(r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redisConn", r)
		c.Next()
	}
}

func main() {
	router := gin.New()

	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	router.Use(APIMiddleware(rClient))

	logrus.SetFormatter(&logrus.JSONFormatter{})
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true, "micro-leaderboard"))

	// Start routes
	router.GET("/health", healthcheck.HealthCheck)
	router.GET("/", app.Root)
	router.GET("/inc/:set/:member", app.Incr)
	router.GET("/inc/:set/:member/:count", app.Incr)
	//router.GET("/member/:set/:member", app.GetRank)
	//router.GET("/board/:set", app.ShowBoard)

	// RUN rabit run
	router.Run() // listen and serve on 0.0.0.0:8080
}
