package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/gin/response"
)

func HelloWorld(c *gin.Context) {
	response.Success(c, "Hello World")
}
