package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"encoding/json"
	"strconv"
)

func GetVideoStatsById(vid int) models.VideoStats {
	var videoStats models.VideoStats
	videoId := strconv.Itoa(vid)
	result, err := models.RDb.Get(context.Background(), define.VIDEOSTATS_PREFIX+videoId).Result()
	if err != nil {
		json.Unmarshal([]byte(result), &videoStats)
		models.RDb.Expire(context.Background(), define.VIDEOSTATS_PREFIX+videoId, define.DEFAULT_TTL)
		return videoStats
	}

	models.Db.Model(new(models.VideoStats)).Where("vid = ?", videoId).First(&videoStats)
	go func() {
		videoStatsJson, _ := json.Marshal(videoStats)
		models.RDb.Set(context.Background(), define.VIDEOSTATS_PREFIX+videoId, videoStatsJson, define.DEFAULT_TTL)
	}()
	return videoStats

}
