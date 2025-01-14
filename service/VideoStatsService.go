package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

func getVideoStatsById(vid int) models.VideoStats {
	var videoStats models.VideoStats

	videoId := strconv.Itoa(vid)
	result, err := models.RDb.Get(context.Background(), define.VIDEOSTATS_PREFIX+videoId).Result()
	if err == nil {
		json.Unmarshal([]byte(result), &videoStats)
		models.RDb.Expire(context.Background(), define.VIDEOSTATS_PREFIX+videoId, define.DEFAULT_TTL)
		return videoStats
	}

	models.Db.Model(new(models.VideoStats)).Where("vid = ?", videoId).Find(&videoStats)
	go func() {
		videoStatsJson, _ := json.Marshal(videoStats)
		models.RDb.Set(context.Background(), define.VIDEOSTATS_PREFIX+videoId, videoStatsJson, define.DEFAULT_TTL)
	}()
	return videoStats

}

func UpdateVideoStats(vid int, field string, icr bool, count int) {
	db := models.Db.Model(new(models.VideoStats)).Where("vid = ?", vid)
	if icr {
		db.Update(field, gorm.Expr(field+"+?", count))
	} else {
		var videoStats models.VideoStats
		db.Find(&videoStats)
		if videoStats.Comment-count > 0 {
			db.Update(field, gorm.Expr(field+"-?", count))
		}
	}
	models.RDb.Del(context.Background(), define.VIDEOSTATS_PREFIX+strconv.Itoa(vid))
}

func UpdateGoodAndBad(vid int, isGood bool) {
	db := models.Db.Model(new(models.VideoStats)).Where("vid = ?", vid)
	if isGood {
		var bad int64
		db.Select("bad").Find(&bad)
		if bad > 0 {
			err := db.Update("bad", gorm.Expr("bad-1")).Error
			if err != nil {
				fmt.Println(err)
			}
		}
		models.Db.Model(new(models.VideoStats)).Where("vid = ?", vid).Update("good", gorm.Expr("good+1"))

	} else {
		var good int64
		db.Select("good").Find(&good)
		if good > 0 {
			db.Update("good", gorm.Expr("good-1"))
		}
		models.Db.Model(new(models.VideoStats)).Where("vid = ?", vid).Update("bad", gorm.Expr("bad+1"))
	}
	models.RDb.Del(context.Background(), define.VIDEOSTATS_PREFIX+strconv.Itoa(vid))
}
