package mid

import (
	"gin-base/pkg/logging"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	logger := logging.GetLogger("Request")
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latencyTime := time.Now().Sub(startTime)
		method := c.Request.Method
		uri := c.Request.RequestURI
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.WithFields(log.Fields{
			"status":       status,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"method":       method,
			"uri":          uri,
		}).Info()
	}
}
