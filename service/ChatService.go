package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetRecentLIst(c *gin.Context) {
	userId, _ := c.Get("userId")
	var chatGetVoList = []models.ChatGetVo{}
	var recentListGetVo = models.RecentListGetVo{}
	offset, _ := strconv.Atoi(c.Query("offset"))
	result, err := models.RDb.ZRange(context.Background(), define.CHAT_ZSET, int64(offset), int64(offset+9)).Result()
	if err != nil || len(result) == 0 {
		recentListGetVo.List = []models.ChatGetVo{}
		response.ResponseOKWithData(c, "获取成功", recentListGetVo)
		return
	}
	var chatList []models.Chat
	models.Db.Model(new(models.Chat)).Where("id in (?)", result).Order("latest_time desc").Find(&chatList)
	for _, v := range chatList {
		var chatGetVo models.ChatGetVo
		chatGetVo.Chat = v
		chatGetVo.User = getUserById(v.UserId)
		chatGetVo.ChatDetail = getChatDetails(v.UserId, userId.(int), 0)
	}

	chatCount, _ := models.RDb.ZCard(context.Background(), define.CHAT_ZSET).Result()
	if offset+10 < int(chatCount) {
		recentListGetVo.More = true
	} else {
		recentListGetVo.More = false
	}
	recentListGetVo.List = chatGetVoList
	response.ResponseOKWithData(c, "获取成功", recentListGetVo)
	return

}
func getChatDetails(uid int, aid int, offset int) models.ChatDetailGetVo {
	var chatDetailGetVo models.ChatDetailGetVo
	key := define.CHAT_DETAILED_ZSET + strconv.Itoa(uid) + ":" + strconv.Itoa(aid)
	chatCount, _ := models.RDb.ZCard(context.Background(), key).Result()
	if offset+20 < int(chatCount) {
		chatDetailGetVo.More = true
	} else {
		chatDetailGetVo.More = false
	}
	chatIdList, _ := models.RDb.ZRange(context.Background(), key, int64(offset), int64(offset+19)).Result()
	if len(chatIdList) == 0 {
		chatDetailGetVo.List = []models.ChatDetail{}
		return chatDetailGetVo
	}
	var chatDetailList []models.ChatDetail
	models.Db.Model(new(models.ChatDetail)).Where("id in (?)", chatIdList).Find(&chatDetailList)
	chatDetailGetVo.List = chatDetailList
	return chatDetailGetVo

}
