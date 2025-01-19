package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
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

func CollectVideo(c *gin.Context) {
	vid, _ := strconv.Atoi(c.PostForm("vid"))
	adds := c.PostForm("adds")
	removes := c.PostForm("removes")
	userId := c.GetInt("userId")
	var addFids []int
	if adds != "" {
		for _, fid := range strings.Split(adds, ",") {
			fidInt, _ := strconv.Atoi(fid)
			addFids = append(addFids, fidInt)
		}
	}
	var removeFids []int
	if removes != "" {
		for _, fid := range strings.Split(removes, ",") {
			fidInt, _ := strconv.Atoi(fid)
			removeFids = append(removeFids, fidInt)
		}
	}

	fids := findFidsOfUser(userId)

	for _, fid := range addFids {
		fmt.Println("fids", fids)
		fmt.Println("fid", fid)
		fmt.Println(isContain(fids, fid))
		if !isContain(fids, fid) {
			response.ResponseFailWithData(c, 403, "无权限", "")
			return
		}
	}
	for _, fid := range removeFids {
		fmt.Println("fids", fids)
		fmt.Println("fid", fid)
		fmt.Println(isContain(fids, fid))
		if !isContain(fids, fid) {
			response.ResponseFailWithData(c, 403, "无权限", "")
			return
		}
	}

	if len(addFids) > 0 {
		addToFavorite(userId, vid, addFids)
	}
	if len(removeFids) > 0 {
		removeFromFavorite(userId, vid, removeFids)
	}

	collectedFids := findFidsOfCollected(vid, fids)
	var isCollect, isCancel bool
	fmt.Println("addFids", addFids)
	fmt.Println("removeFids", removeFids)
	fmt.Println("collectedFids", collectedFids)
	if len(addFids) > 0 && len(removeFids) == 0 {
		fmt.Println("1")
		isCollect = true
	}
	if len(addFids) == 0 && len(removeFids) > 0 && len(removeFids) == len(collectedFids) && isContainAll(collectedFids, removeFids) {
		fmt.Println("2")
		isCancel = true
	}

	if isCollect {
		collectOrCancel(userId, vid, true)
	} else if isCancel {
		collectOrCancel(userId, vid, false)
	}
	response.ResponseOK(c)
	return

}

func GetCollectedFids(c *gin.Context) {
	vid, _ := strconv.Atoi(c.Query("vid"))
	userId := c.GetInt("userId")
	fids := findFidsOfUser(userId)
	var collectedFids []int
	models.Db.Model(new(models.FavoriteVideo)).
		Select("fid").
		Where("vid = ?", vid).
		Where("fid in ?", fids).
		Where("is_remove = ?", 0).
		Find(&collectedFids)
	response.ResponseOKWithData(c, "", collectedFids)
	return
}

func AddFavorite(c *gin.Context) {
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	visible, _ := strconv.Atoi(c.PostForm("visible"))
	userId := c.GetInt("userId")

	var favorite models.Favorite
	favorite.Uid = userId
	favorite.Title = title
	favorite.Visible = visible
	favorite.Description = desc
	if err := models.Db.Create(&favorite).Error; err != nil {
		response.ResponseFailWithData(c, 403, "创建失败", "")
		return
	}
	models.RDb.Del(context.Background(), define.FAVORITE_PREFIX+strconv.Itoa(userId))
	response.ResponseOKWithData(c, "", favorite)
	return

}

func findFidsOfCollected(vid int, fids []int) []int {
	if len(fids) == 0 {
		return []int{}
	}
	var existFids []int
	models.Db.Model(new(models.FavoriteVideo)).
		Select("fid").
		Where("vid = ?", vid).
		Where("fid in ?", fids).
		Find(&existFids)
	return existFids
}

func addToFavorite(uid int, vid int, fids []int) {
	var existFids []int

	models.Db.Model(new(models.FavoriteVideo)).
		Select("fid").
		Where("vid = ?", vid).
		Where("fid in ?", fids).
		Find(&existFids)

	models.Db.Model(new(models.FavoriteVideo)).
		Where("vid = ?", vid).
		Where("fid in ?", fids).Update("is_remove", 0)

	finalIds := diffOfInt(fids, existFids)

	for _, fid := range finalIds {
		models.Db.Create(&models.FavoriteVideo{
			Fid:  fid,
			Vid:  vid,
			Time: models.MyTime(time.Now()),
		})
	}

	models.Db.Model(new(models.Favorite)).
		Where("fid in ?", fids).
		Update("count", gorm.Expr("count + ?", 1))

	models.RDb.ZAdd(context.Background(), define.USER_FAVORITE_VIDEO, &redis.Z{
		Member: vid,
		Score:  float64(time.Now().Unix()),
	})
	models.RDb.Del(context.Background(), define.USER_FAVORITES+strconv.Itoa(uid))

}

func removeFromFavorite(uid int, vid int, fids []int) {
	models.Db.Model(new(models.FavoriteVideo)).
		Where("vid = ?", vid).
		Where("fid in ?", fids).Update("is_remove", 1)

	models.Db.Model(new(models.Favorite)).
		Where("fid in ?", fids).
		Where("count > ?", 0).
		Update("count", gorm.Expr("count - ?", 1))

	models.RDb.ZRem(context.Background(), define.USER_FAVORITE_VIDEO, vid)
	models.RDb.Del(context.Background(), define.USER_FAVORITES+strconv.Itoa(uid))
}

func diffOfInt(a, b []int) []int {
	var diff []int
	for _, v := range a {
		if !isContain(b, v) {
			diff = append(diff, v)
		}
	}
	return diff

}

func isContainAll(a, b []int) bool {
	for _, v := range b {
		if !isContain(a, v) {
			return false
		}
	}
	return true
}
func isContain(a []int, x int) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func findFidsOfUser(userId int) []int {
	var fids []int
	favorites := getFavorites(userId, true)
	for _, v := range favorites {
		fids = append(fids, v.Fid)
	}
	return fids
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
		Order("fid desc").Find(&favoriteList)
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
