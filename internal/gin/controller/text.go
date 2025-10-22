package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/internal/gin/response"
)

//文本消息处理器

func HandleMsg(c *gin.Context) {
	response.Success(c, "success")
}

func TestError(c *gin.Context) {
	response.Fail(c, http.StatusBadRequest, "接口异常", "线程池耗尽")
}
