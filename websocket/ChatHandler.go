package websocketServer

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func sendWhisper(imMessage models.ImMessage) {
	fmt.Println(imMessage)
	fromId := imMessage.UserId
	messageMap := imMessage.Message
	content := messageMap["content"].(string)
	toId := int(messageMap["anotherId"].(float64))
	chatDetail := models.ChatDetail{
		UserId:    fromId,
		AnotherId: toId,
		Content:   content,
		Time:      models.MyTime(time.Now()),
	}
	models.Db.Model(new(models.ChatDetail)).Create(&chatDetail)
	models.RDb.ZAdd(context.Background(), define.CHAT_DETAILED_ZSET+strconv.Itoa(fromId)+":"+strconv.Itoa(toId), &redis.Z{
		Member: chatDetail.Id,
		Score:  float64(time.Now().Unix()),
	})
	models.RDb.ZAdd(context.Background(), define.CHAT_DETAILED_ZSET+strconv.Itoa(toId)+":"+strconv.Itoa(fromId), &redis.Z{
		Member: chatDetail.Id,
		Score:  float64(time.Now().Unix()),
	})

}
