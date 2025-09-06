package middleware

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stoppieboy/rate-limiter-server/internal/metrics"
)

var (
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

type Data struct {
	Allowed          int64
	Remaining_tokens int64
	MS               int64
	Server_time      int64
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("token", "12345")
		// redis script calling
		script, err := os.ReadFile("token_bucket.lua")
		tokenBucketScript := redis.NewScript(string(script))
		keys := []string{"bucket"}
		// capacity, refill_rate, requested, now_millis
		args := []interface{}{100, 2, 50}
		scriptStart := time.Now()
		data, err := tokenBucketScript.Run(ctx, Rdb, keys, args).Result()
		metrics.RedisOpDuration.WithLabelValues("lua_token_bucket").Observe(time.Since(scriptStart).Seconds())
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
		path := c.Request.URL.String()
		if val.Allowed == 1 {
			c.Next()
			metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status()))
			metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(t).Seconds())
		} else {
			metrics.RateLimitedTotal.WithLabelValues(path).Inc()
			metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, path, "429")
			c.AbortWithStatusJSON(429, gin.H{"message": "Rate Limit Exceeded"})
			metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(t).Seconds())
		}

		latency := time.Since(t)
		log.Print("Latency: ", latency)
	}
}
