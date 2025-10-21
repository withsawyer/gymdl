package middleware

import (
	"log"
	"os"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/config"
)

func LoggerMiddleware(logConfig *config.LogConfig) gin.HandlerFunc {
	//todo 日志级别控制 以及模式3
	if logConfig.Mode != 1 {
		f, err := os.OpenFile(logConfig.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		logger := log.New(f, "", log.LstdFlags)
		
		return func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path
			method := c.Request.Method
			c.Next()
			duration := time.Since(start)
			statusCode := c.Writer.Status()
			logger.Printf("[%s] %s %s %d %s",
				time.Now().Format("2006-01-02 15:04:05"),
				method,
				path,
				statusCode,
				duration,
			)
		}
	} else {
		return func(context *gin.Context) {}
	}
}
