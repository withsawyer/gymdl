package middleware

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/utils"
)

func LoggerMiddleware() gin.HandlerFunc {
	log := utils.Logger()
	return func(c *gin.Context) {
		start := time.Now() // 记录请求开始时间
		
		c.Next() // 处理请求
		
		// 请求处理完成后，记录日志
		latency := time.Since(start)    // 计算耗时
		statusCode := c.Writer.Status() // HTTP状态码
		clientIP := c.ClientIP()        // 客户端IP
		method := c.Request.Method      // 请求方法
		path := c.Request.URL.Path      // 请求路径
		bodySize := c.Writer.Size()     // 响应体大小
		
		log.Infof("%3d | %13v | %15s | %-7s %s | size: %d",
			statusCode, latency, clientIP, method, path, bodySize)
	}
}
