package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/internal/gin/middleware"
)

// SetupRouter 路由注册
func SetupRouter(c *config.Config) *gin.Engine {
	engine := gin.New()
	//中间件注册(扩展) 常见中间件有日志、鉴权、限流、跨域等
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.LoggerMiddleware())
	apiGroup := engine.Group("/api")
	//注册文本处理器路由
	RegisterTextRoutes(apiGroup)
	//注册指令处理器路由
	RegisterCommandRoutes(apiGroup)
	return engine
}
