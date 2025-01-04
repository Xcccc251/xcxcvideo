package service

import (
	"XcxcVideo/common/response"
	"github.com/gin-gonic/gin"
)

func SearchHotList(c *gin.Context) {
	response.ResponseOKWithData(c, "获取成功", nil)
}
