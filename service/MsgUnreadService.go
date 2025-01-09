package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
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
