package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/gin/controller"
)

func RegisterTextRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/text")
	group.GET("/", controller.HandleMsg)
}
