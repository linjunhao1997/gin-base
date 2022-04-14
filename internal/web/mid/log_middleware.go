package mid

import (
	"fmt"
	"gin-base/pkg/logging"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

var logger = logging.GetLogger("Request")

func LogMiddleware() gin.HandlerFunc {
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
			"latency_time": fmt.Sprintf("%v", latencyTime),
			"client_ip":    clientIP,
			"method":       method,
			"uri":          uri,
		}).Info()
	}
}
