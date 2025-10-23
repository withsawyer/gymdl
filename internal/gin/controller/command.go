package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/gin/response"
)

// 指令处理器

func HandleCommand(c *gin.Context) {
	response.Success(c, gin.H{
		"hahah": "sd",
		"data":  "dsd",
	})
}
