package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	websocketServer "XcxcVideo/websocket"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
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
		chatGetVoList = append(chatGetVoList, chatGetVo)
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
func CreateChat(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	userId, _ := c.Get("userId")
	var count int64
	db := models.Db.Model(new(models.Chat)).Where("user_id = ? and another_id = ?", uid, userId.(int))
	db.Count(&count)
	var chat models.Chat
	if count != 0 {
		db.Find(&chat)
		if chat.IsDeleted == 1 {
			db.Update("is_deleted", 0)
			db.Update("latest_time", models.MyTime(time.Now()))
			models.RDb.ZAdd(context.Background(), define.CHAT_ZSET, &redis.Z{
				Member: chat.Id,
				Score:  float64(time.Now().Unix()),
			})
			var rmap = make(map[string]interface{})
			rmap["chat"] = chat
			sw := sync.WaitGroup{}
			sw.Add(2)
			go func() {
				rmap["user"] = getUserById(chat.UserId)
				sw.Done()
			}()
			go func() {
				rmap["detail"] = getChatDetails(uid, userId.(int), 0)
				sw.Done()
			}()
			sw.Wait()
			rmap["msg"] = "新创建"
			response.ResponseOKWithData(c, "获取成功", rmap)
			return

		} else {
			response.ResponseOKWithData(c, "获取成功", gin.H{
				"msg": "已存在",
			})
			return
		}

	} else {

		var userCount int64
		models.Db.Model(new(models.User)).Where("id = ?", uid).Count(&userCount)
		if userCount == 0 {
			response.ResponseFailWithData(c, 500, "用户不存在", nil)
			return
		}
		newChat := models.Chat{
			UserId:     uid,
			AnotherId:  userId.(int),
			LatestTime: models.MyTime(time.Now()),
		}
		models.Db.Model(new(models.Chat)).Create(&newChat)
		models.RDb.ZAdd(context.Background(), define.CHAT_ZSET, &redis.Z{
			Member: newChat.Id,
			Score:  float64(time.Now().Unix()),
		})
		var rmap = make(map[string]interface{})
		rmap["chat"] = newChat
		sw := sync.WaitGroup{}
		sw.Add(2)
		go func() {
			rmap["user"] = getUserById(newChat.UserId)
			sw.Done()
		}()
		go func() {
			rmap["detail"] = getChatDetails(uid, userId.(int), 0)
			sw.Done()
		}()
		sw.Wait()
		rmap["msg"] = "新创建"
		response.ResponseOKWithData(c, "获取成功", rmap)
		return
	}
}

func UpdateWhisperOnline(c *gin.Context) {
	from, _ := strconv.Atoi(c.Query("from"))
	userId, _ := c.Get("userId")
	key := define.WHISPER_KEY + strconv.Itoa(userId.(int)) + ":" + strconv.Itoa(from)
	models.RDb.Set(context.Background(), key, 1, 0)
	var chat models.Chat
	db := models.Db.Model(new(models.Chat)).Where("user_id = ? and another_id = ?", from, userId.(int))
	db.Find(&chat)
	if chat.Unread > 0 {
		messageMap := map[string]interface{}{}
		messageMap["type"] = "已读"
		messageMap["id"] = chat.Id
		messageMap["count"] = chat.Unread
		websocketServer.SendMessage(userId.(int), "whisper", messageMap)
		subtractWhisper(userId.(int), chat.Unread)
	}
	db.Update("unread", 0)
	response.ResponseOK(c)
	return
}

func UpdateWhisperOutline(c *gin.Context) {
	//放开该接口，无需验证token
	from, _ := strconv.Atoi(c.Query("from"))
	to, _ := strconv.Atoi(c.Query("to"))
	key := define.WHISPER_KEY + strconv.Itoa(to) + ":" + strconv.Itoa(from)
	models.RDb.Del(context.Background(), key)
	response.ResponseOK(c)
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

//func updateChat(from int, to int) {
//	key := define.WHISPER_KEY + strconv.Itoa(to) + ":" + strconv.Itoa(from)
//	result, _ := models.RDb.Exists(context.Background(), key).Result()
//
//}
