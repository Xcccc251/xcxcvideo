package websocketServer

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/commonService"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
	"time"
)

func SendWhisper(imMessage models.ImMessage) {
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
	online := commonService.UpdateChat(fromId, toId)
	imMessageMap := map[string]interface{}{
		"type":   "接收",
		"online": online,
		"detail": chatDetail,
	}
	sw := sync.WaitGroup{}
	sw.Add(2)
	go func() {
		var chat models.Chat
		models.Db.Model(new(models.Chat)).Where("user_id = ? and another_id = ?", fromId, toId).Find(&chat)
		imMessageMap["chat"] = chat
		sw.Done()
	}()
	go func() {
		imMessageMap["user"] = commonService.GetUserById(fromId)
		sw.Done()
	}()
	sw.Wait()
	fmt.Println("发送消息：", toId, imMessageMap)
	SendMessage(toId, "whisper", imMessageMap)
	fmt.Println("发送消息：", fromId, imMessageMap)
	SendMessage(fromId, "whisper", imMessageMap)

}
