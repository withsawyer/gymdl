package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/gin/response"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
)

// bodyWriter 包装 gin.ResponseWriter 用于捕获响应体
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinLoggerMiddleware 日志中间件
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 包装 ResponseWriter
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw
		
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path += "?" + raw
		}
		
		method := c.Request.Method
		clientIP := c.ClientIP()
		
		c.Next()
		
		cost := time.Since(start)
		
		// 默认日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("latency", cost.String()),
		}
		
		// 尝试解析响应体 JSON
		var resp response.Response
		bodyBytes := bw.body.Bytes()
		if err := json.Unmarshal(bodyBytes, &resp); err == nil {
			fields = append(fields,
				zap.Int("code", resp.Code),
				zap.String("message", resp.Message),
				zap.Strings("errors", resp.Errors),
				//zap.Any("data", resp.Data),
			)
		}
		
		// 根据业务 Code 判断日志等级
		code := resp.Code
		logger := utils.Logger()
		switch {
		case code >= 500:
			logger.Error("[GIN]", fields...)
		case code >= 400:
			logger.Warn("[GIN]", fields...)
		default:
			logger.Info("[GIN]", fields...)
		}
	}
}

// GinRecoveryMiddleware 优化版
func GinRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				req := c.Request
				requestInfo := req.Method + " " + req.URL.Path + " " + req.Proto
				logger := utils.Logger()
				
				// 检查网络断开错误
				if isBrokenPipeErr(rec) {
					logger.Warn("[BROKEN PIPE]",
						zap.Any("error", rec),
						zap.String("request", requestInfo),
					)
					c.Abort()
					return
				}
				
				// 其他 panic，返回 Fail 响应
				logger.Error("[PANIC RECOVER]",
					zap.Any("error", rec),
					zap.String("request", requestInfo),
				)
				response.Fail(c, 500, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// isBrokenPipeErr 检查网络断开
func isBrokenPipeErr(rec any) bool {
	err, ok := rec.(error)
	if !ok || err == nil {
		return false
	}
	
	var netErr net.Error
	if errors.As(err, &netErr) {
		msg := err.Error()
		return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
	}
	return false
}
