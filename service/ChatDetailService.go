package service

import (
	"XcxcVideo/common/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetMoreChatDetail(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	userId, _ := c.Get("userId")
	chatDetails := getChatDetails(uid, userId.(int), offset)
	response.ResponseOKWithData(c, "", chatDetails)
	return
}
