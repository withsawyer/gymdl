package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/controller"
)

func RegisterCommandRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/command")
	group.GET("/", controller.HandleCommand)
}
