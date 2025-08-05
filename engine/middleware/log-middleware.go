package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/controller"
	"go-quantus-service/src/entities"
	"io/ioutil"
	"strconv"
	"time"
)

func RequestLogger(logController controller.LogControllerinterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read body
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Proses request
		c.Next()

		// After request
		status := c.Writer.Status()
		headers, _ := json.Marshal(c.Request.Header)

		log := entities.LogEntry{
			IPAddress: c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Headers:   string(headers),
			Body:      string(bodyBytes),
			Response:  strconv.Itoa(status),
			Status:    status,
			CreatedAt: start,
		}

		// Ambil Redis dari controller
		_ = logController.GetDependencies().Redis.PushLogToQueue("log_queue", log)
	}
}
