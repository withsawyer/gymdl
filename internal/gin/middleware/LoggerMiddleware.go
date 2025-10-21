package middleware

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/utils"
)

// GinLoggerMiddleware 是一个 Gin 中间件，用于记录请求日志
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 执行请求
		c.Next()

		cost := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errs := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		utils.Logger().Infow("GIN",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", fmt.Sprintf("%v", cost),
			"errors", errs,
		)
	}
}

// GinRecoveryMiddleware 捕获 panic 并使用 zap 打印日志
func GinRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(error); ok {
					if opErrMsg := ne.Error(); bytes.Contains([]byte(opErrMsg), []byte("broken pipe")) ||
						bytes.Contains([]byte(opErrMsg), []byte("connection reset by peer")) {
						brokenPipe = true
					}
				}

				httpRequest, _ := dumpRequest(c)
				if brokenPipe {
					utils.Logger().Errorw("BROKEN PIPE", "error", err, "request", string(httpRequest))
					c.Abort()
					return
				}

				utils.Logger().Errorw("PANIC RECOVER",
					"error", err,
					"request", string(httpRequest),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// dumpRequest 简化版请求信息输出
func dumpRequest(c *gin.Context) (string, error) {
	req := c.Request
	return fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, req.Proto), nil
}
