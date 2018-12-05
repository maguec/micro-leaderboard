package main

import (
	"fmt"
	"os"
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
	var redisHost string
	var redisPort string
	var listenPort string

	if len(os.Getenv("REDIS_HOST")) > 0 {
		redisHost = os.Getenv("REDIS_HOST")
	} else {
		redisHost = "localhost"
	}

	if len(os.Getenv("REDIS_PORT")) > 0 {
		redisPort = os.Getenv("REDIS_PORT")
	} else {
		redisPort = "6379"
	}

	if len(os.Getenv("LISTEN_PORT")) > 0 {
		listenPort = os.Getenv("LISTEN_PORT")
	} else {
		listenPort = "8080"
	}

	rClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
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
	router.GET("/member/:set/:member", app.GetRank)
	router.GET("/board/:set", app.ShowBoard)
	router.GET("/board/:set/:count", app.ShowBoard)

	// RUN rabit run
	router.Run(fmt.Sprintf(":%s", listenPort)) // listen and serve on 0.0.0.0:8080
}
