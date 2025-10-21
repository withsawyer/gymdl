package controller

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

//指令处理器

func HandleCommand(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "指令处理器"})
}
