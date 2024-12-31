package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetAllFavoritesForUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	userId, _ := c.Get("userId")
	var favoriteList []models.Favorite
	if uid == userId {
		favoriteList = getFavorites(uid, true)
	} else {
		favoriteList = getFavorites(uid, false)
	}
	response.ResponseOKWithData(c, "", favoriteList)
	return

}
func getFavorites(uid int, isOwner bool) []models.Favorite {
	key := define.FAVORITE_PREFIX + strconv.Itoa(uid)
	result, err := models.RDb.Get(context.Background(), key).Result()
	if err == nil {
		var favoriteList []models.Favorite
		var resultList []models.Favorite
		json.Unmarshal([]byte(result), &favoriteList)
		if !isOwner {
			for _, v := range favoriteList {
				if v.Visible == 1 {
					resultList = append(resultList, v)
				}
			}
			return resultList
		}
		return favoriteList
	}

	var favoriteList []models.Favorite
	models.Db.Model(new(models.Favorite)).
		Where("uid = ?", uid).
		Order("id desc")
	models.Db.Find(&favoriteList)
	for _, v := range favoriteList {
		if v.Cover == "" {
			set := redisUtil.GetSet(define.FAVORITE_VIDEO_PREFIX + strconv.Itoa(v.Fid))
			if len(set) > 0 {
				vid, _ := strconv.Atoi(set[0])
				var video models.Video
				models.Db.Model(new(models.Video)).Where("id = ?", vid).Find(&video)
				v.Cover = video.CoverUrl
			}
		}

	}
	resultList := favoriteList
	favoriteListJson, _ := json.Marshal(favoriteList)
	models.RDb.Set(context.Background(), key, favoriteListJson, define.DEFAULT_TTL)
	if !isOwner {
		for _, v := range favoriteList {
			if v.Visible == 1 {
				resultList = append(resultList, v)
			}
		}
		return resultList
	}
	return favoriteList

}
