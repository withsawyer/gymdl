package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//文本消息处理器

func HandleMsg(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "消息处理器"})
}

func TestError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"data": "服务器内部异常"})
}
