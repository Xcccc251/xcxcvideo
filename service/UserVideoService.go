package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func UserPlayVideo(c *gin.Context) {
	vid, _ := strconv.Atoi(c.PostForm("vid"))
	userId, _ := c.Get("userId")
	userVideo := updatePlay(vid, userId.(int))
	response.ResponseOKWithData(c, "", userVideo)
	return
}
func updatePlay(vid int, uid int) models.UserVideo {
	dbUserVideo := models.UserVideo{}
	var count int64
	db := models.Db.Model(new(models.UserVideo)).Where("vid = ? and uid = ?", vid, uid)
	db.Find(&dbUserVideo)
	db.Count(&count)
	if count == 0 {
		dbUserVideo.Uid = uid
		dbUserVideo.Vid = vid
		dbUserVideo.Play = 1
		dbUserVideo.PlayTime = models.MyTime(time.Now())
		db.Create(&dbUserVideo)
	} else if time.Now().Unix()-time.Time(dbUserVideo.PlayTime).Unix() < 30 {
		fmt.Println("间隔小于30s")
		return dbUserVideo
	} else {
		db.Update("play", dbUserVideo.Play+1)
		db.Update("play_time", models.MyTime(time.Now()))
	}
	//todo 消息队列
	go func() {
		models.RDb.ZAdd(context.Background(), define.USER_VIDEO_HISTORY+strconv.Itoa(vid), &redis.Z{
			Member: vid,
			Score:  float64(time.Now().Unix()),
		})
		updateVideoStats(vid, "play", true, 1)
	}()
	return dbUserVideo
}
