package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
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

func DeleteChatDetail(c *gin.Context) {
	chatDetailId, _ := strconv.Atoi(c.PostForm("id"))
	uid, _ := c.Get("userId")
	var count int64
	var chatDetail models.ChatDetail
	db := models.Db.Model(new(models.ChatDetail)).Where("id=?", chatDetailId)
	db.Count(&count)
	if count == 0 {
		response.ResponseFail(c)
		return
	}
	db.Find(&chatDetail)
	if chatDetail.UserId == uid {
		db.Update("user_del", 1)
		key := define.CHAT_DETAILED_ZSET + strconv.Itoa(chatDetail.AnotherId) + ":" + strconv.Itoa(uid.(int))
		models.RDb.ZRem(context.Background(), key, chatDetailId)
		response.ResponseOK(c)
		return
	} else {
		db.Update("user_del", 1)
		key := define.CHAT_DETAILED_ZSET + strconv.Itoa(chatDetail.UserId) + ":" + strconv.Itoa(uid.(int))
		models.RDb.ZRem(context.Background(), key, chatDetailId)
		response.ResponseOK(c)
		return
	}
	response.ResponseFail(c)
	return

}
