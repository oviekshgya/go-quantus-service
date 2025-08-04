package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/entities"
	rds "go-quantus-service/src/redis"
	"io/ioutil"
	"strconv"
	"time"
)

func RequestLogger(redisClient *rds.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read body (and replace so it can be read again)
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Continue processing
		c.Next()

		// Capture response body if needed (optional: or use middleware like gin's ResponseRecorder)
		status := c.Writer.Status()

		headers, _ := json.Marshal(c.Request.Header)

		log := entities.LogEntry{
			IPAddress: c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Headers:   string(headers),
			Body:      string(bodyBytes),
			Response:  strconv.Itoa(c.Writer.Status()),
			Status:    status,
			CreatedAt: start,
		}

		_ = redisClient.PushLogToQueue("log_queue", log)

	}
}
