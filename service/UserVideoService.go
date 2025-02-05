package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	websocketServer "XcxcVideo/handler"
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

func LikeOrDisLikeVideo(c *gin.Context) {
	vid, _ := strconv.Atoi(c.PostForm("vid"))
	isLove, _ := strconv.ParseBool(c.PostForm("isLove"))
	isSet, _ := strconv.ParseBool(c.PostForm("isSet"))
	userId, _ := c.Get("userId")
	key := define.LOVE_VIDEO + strconv.Itoa(userId.(int))
	var userVideo models.UserVideo
	db := models.Db.Model(new(models.UserVideo)).Where("vid = ? and uid = ?", vid, userId)
	db.Find(&userVideo)

	if isLove && isSet {
		if userVideo.Love == 1 {
			response.ResponseOKWithData(c, "已点赞", userVideo)
			return
		}
		userVideo.Love = 1
		db.Update("love", 1)
		db.Update("love_time", models.MyTime(time.Now()))

		if userVideo.Unlove == 1 {
			userVideo.Unlove = 0
			db.Update("unlove", 0)
			go func() {
				UpdateGoodAndBad(vid, true)
			}()
		} else {
			go func() {
				UpdateVideoStats(vid, "good", true, 1)
			}()
		}
		models.RDb.ZAdd(context.Background(), key, &redis.Z{
			Member: vid,
			Score:  float64(time.Now().Unix()),
		})

		go func() {
			var video models.Video
			models.Db.Model(new(models.Video)).Where("id = ?", vid).Find(&video)
			models.RDb.ZAdd(context.Background(), define.BELOVED_VIDEO_SET+strconv.Itoa(video.Uid), &redis.Z{
				Member: vid,
				Score:  float64(time.Now().Unix()),
			})
			addOneUnreadMsg(video.Uid, "love")
			var messageMap = map[string]string{
				"type": "接收",
			}
			websocketServer.SendMessage(video.Uid, "love", messageMap)
		}()

	} else if isLove {
		if userVideo.Love == 0 {
			response.ResponseOKWithData(c, "已取消", userVideo)
			return
		}
		userVideo.Love = 0
		db.Update("love", 0)
		models.RDb.ZRem(context.Background(), key, vid)
		go func() {
			UpdateVideoStats(vid, "good", false, 1)
		}()
	} else if isSet {
		if userVideo.Unlove == 1 {
			response.ResponseOKWithData(c, "已点踩", userVideo)
			return
		}
		userVideo.Unlove = 1
		db.Update("unlove", 1)
		if userVideo.Love == 1 {
			userVideo.Love = 0
			db.Update("love", 0)
			go func() {
				UpdateGoodAndBad(vid, false)
			}()
		} else {
			go func() {
				UpdateVideoStats(vid, "bad", true, 1)
			}()
		}
	} else {
		if userVideo.Unlove == 0 {
			response.ResponseOKWithData(c, "已取消", userVideo)
			return
		}
		db.Update("unlove", 0)
		userVideo.Unlove = 0
		go func() {
			UpdateVideoStats(vid, "bad", false, 1)
		}()
	}
	response.ResponseOKWithData(c, "点赞成功", userVideo)
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
	go func() {
		models.RDb.ZAdd(context.Background(), define.USER_VIDEO_HISTORY+strconv.Itoa(vid), &redis.Z{
			Member: vid,
			Score:  float64(time.Now().Unix()),
		})
		UpdateVideoStats(vid, "play", true, 1)
	}()
	return dbUserVideo
}

func collectOrCancel(uid int, vid int, isCollect bool) {
	db := models.Db.Model(new(models.UserVideo)).
		Where("vid = ? and uid = ?", vid, uid)
	if isCollect {
		db.Update("collect", 1)
	} else {
		db.Update("collect", 0)
	}

	UpdateVideoStats(vid, "collect", isCollect, 1)
}
