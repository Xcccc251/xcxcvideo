package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	websocketServer "XcxcVideo/handler"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetMsgUnread(c *gin.Context) {
	userId, _ := c.Get("userId")
	uidStr := strconv.Itoa(userId.(int))
	var msgUnread models.MsgUnread
	resultJson, err := models.RDb.Get(context.Background(), define.MSG_UNREAD+uidStr).Result()
	if err == redis.Nil {
		models.Db.Model(new(models.MsgUnread)).Where("uid = ?", userId).Find(&msgUnread)
		msgUnreadJson, _ := json.Marshal(msgUnread)
		models.RDb.Set(context.Background(), define.MSG_UNREAD+uidStr, msgUnreadJson, define.MSG_UNREAD_TTL)
		response.ResponseOKWithData(c, "", msgUnread)
		return
	}
	json.Unmarshal([]byte(resultJson), &msgUnread)
	response.ResponseOKWithData(c, "", msgUnread)
	return
}

func addOneUnreadMsg(userId int, column string) {
	models.Db.Model(new(models.MsgUnread)).Where("uid = ?", userId).Update(column, gorm.Expr(column+"+1"))
	models.RDb.Del(context.Background(), define.MSG_UNREAD+strconv.Itoa(userId))
}

func ClearUnreadMsg(c *gin.Context) {
	userId, _ := c.Get("userId")
	column := c.PostForm("column")
	var count int64
	models.Db.Model(new(models.MsgUnread)).Select(column).Where("uid = ?", userId).Count(&count)
	if count == 0 {
		response.ResponseOK(c)
		return
	}
	imresponse := websocketServer.ImResponse{
		Type: column,
		Time: models.MyTime(time.Now()),
		Data: map[string]interface{}{
			"type": "全部已读",
		},
	}
	websocketServer.SendMessage(userId.(int), imresponse.Type, imresponse.Data)
	models.Db.Model(new(models.MsgUnread)).Where("uid = ?", userId).Update(column, 0)
	models.RDb.Del(context.Background(), define.MSG_UNREAD+strconv.Itoa(userId.(int)))
	//todo 私聊
	response.ResponseOK(c)
	return
}

func subtractWhisper(userId int, count int) {
	db := models.Db.Model(new(models.MsgUnread)).Where("uid = ?", userId)
	var whisperCount int64
	db.Select("whisper").Find(&whisperCount)
	if whisperCount < int64(count) {
		db.Update("whisper", 0)
	} else {
		db.Update("whisper", gorm.Expr("whisper-?", count))
	}
	models.RDb.Del(context.Background(), define.MSG_UNREAD+strconv.Itoa(userId))
}
