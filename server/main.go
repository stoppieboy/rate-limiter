package main

import (
	"context"
	// "encoding/json"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

type Data struct {
	Allowed int64
	Remaining_tokens int64
	MS int64
	Server_time int64
}

func RateLimiter() gin.HandlerFunc {
	return func (c *gin.Context){
		t := time.Now()
		c.Set("token", "12345")
		// redis script calling
		script, err := os.ReadFile("token_bucket.lua")
		tokenBucketScript := redis.NewScript(string(script))
		keys := []string{"bucket"}
		// capacity, refill_rate, requested, now_millis
		args := []interface{}{100, 2, 50}
		data, err := tokenBucketScript.Run(ctx, rdb, keys, args).Result()
		// err = rdb.Set(ctx, "car", "Lambo", 0).Err()
		if err != nil {
			panic(err)
		}
		// allowed, remaining tokens, ms until next token, server time
		log.Print(data)
		slice := data.([]interface{})
		var val Data

		val.Allowed = slice[0].(int64)
		val.Remaining_tokens = slice[1].(int64)
		val.MS = slice[2].(int64)
		val.Server_time = slice[3].(int64)

		log.Print(val)

		// byteData, _ := json.Marshal(data)
		// receivedData := Data{}
		// json.Unmarshal(byteData, &receivedData)

		// log.Printf("Allowed: %v, remaining tokens: %v, MS until next token: %v, server time: %v\n",receivedData.Allowed, receivedData.Remaining_tokens, receivedData.MS, receivedData.Server_time)
		if val.Allowed == 1{
			c.Next()
		} else {
			c.AbortWithStatusJSON(429, gin.H{"message": "Rate Limit Exceeded"})
		}


		latency := time.Since(t)
		log.Print("Latency: ",latency)
	}
}

func main() {
	defer rdb.Close()
	router := gin.Default()
	router.Use(RateLimiter())
	router.GET("/ping", func( c *gin.Context) {
		token := c.MustGet("token").(string)
		for i:=0;i<100000000;i++{}
		c.JSON(200, gin.H{
			"message": "pong",
			"token": token,
		})
	})

	router.Run(":3000")
}