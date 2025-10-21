package controller

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

//文本消息处理器

func HandleMsg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "消息处理器"})
}
