package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stoppieboy/rate-limiter-server/internal/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
	defer middleware.Rdb.Close()
	router := gin.Default()
	// router.Use(middleware.RateLimiter())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/ping",  middleware.RateLimiter(), func( c *gin.Context) {
		token := c.MustGet("token").(string)
		// for i:=0;i<100000000;i++{}
		c.JSON(200, gin.H{
			"message": "pong",
			"token": token,
		})
	})



	router.Run(":3000")
}