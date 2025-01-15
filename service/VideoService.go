package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/es"
	"XcxcVideo/common/helper"
	"XcxcVideo/common/minIO"
	"XcxcVideo/common/models"
	"XcxcVideo/common/oss"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetVideoById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("vid"))
	video := getVideoById(id)
	if video.Video.Status != 1 {
		response.ResponseFailWithData(c, 404, "视频不存在", "")
		return
	}
	response.ResponseOKWithData(c, "", video)
	return
}

func ChangeVideoStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("vid"))
	status, _ := strconv.Atoi(c.PostForm("status"))
	userId, _ := c.Get("userId")
	if status == 1 {
		var video models.VideoVo
		db := models.Db.Model(new(models.VideoVo)).Where("id = ?", id)
		var count int64
		db.Count(&count)
		db.Find(&video)
		if count == 0 {
			response.ResponseFailWithData(c, 404, "视频不见了", "")
			return
		}
		err := db.Update("status", status).Error
		if err != nil {
			response.ResponseFailWithData(c, 500, "服务器错误", "")
			return
		}
		lastStatus := video.Status
		video.Status = status
		es.AddSearchVideo(video)
		models.RDb.SAdd(context.Background(), define.USER_VIDEO_UPLOAD+strconv.Itoa(video.Uid), video.Vid)
		models.RDb.SRem(context.Background(), define.VIDEO_STATUS+strconv.Itoa(lastStatus), video.Vid)
		models.RDb.SAdd(context.Background(), define.VIDEO_STATUS+strconv.Itoa(status), video.Vid)
		models.RDb.Del(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(video.Vid))
		response.ResponseOK(c)
		return

	} else if status == 2 {
		var video models.VideoVo
		db := models.Db.Model(new(models.VideoVo)).Where("id = ?", id)
		var count int64
		db.Count(&count)
		db.Find(&video)
		if count == 0 {
			response.ResponseFailWithData(c, 404, "视频不见了", "")
			return
		}
		err := db.Update("status", status).Error
		if err != nil {
			response.ResponseFailWithData(c, 500, "服务器错误", "")
			return
		}
		lastStatus := video.Status
		video.Status = status
		models.RDb.SRem(context.Background(), define.USER_VIDEO_UPLOAD+strconv.Itoa(video.Uid), video.Vid)
		models.RDb.SRem(context.Background(), define.VIDEO_STATUS+strconv.Itoa(lastStatus), video.Vid)
		models.RDb.SAdd(context.Background(), define.VIDEO_STATUS+strconv.Itoa(status), video.Vid)
		models.RDb.Del(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(video.Vid))
		response.ResponseOK(c)
		return
	} else {
		var video models.VideoVo
		db := models.Db.Model(new(models.VideoVo)).Where("id = ?", id)
		var count int64
		db.Count(&count)
		db.Find(&video)
		if count == 0 {
			response.ResponseFailWithData(c, 404, "视频不见了", "")
			return
		}
		if userId.(int) == video.Uid || IsAdmin(userId.(int)) {
			videoUrl := video.VideoUrl
			videoUrl = strings.Split(videoUrl, "xcxcaudio/")[1]
			coverUrl := video.CoverUrl
			coverUrl = strings.Split(coverUrl, "aliyuncs.com/")[1]
			err := db.Update("status", 3).Error
			db.Delete(new(models.Video))
			if err != nil {
				response.ResponseFailWithData(c, 500, "服务器错误", "")
				return
			}
			es.UpdateVideoStatus(video.Vid, status)
			lastStatus := video.Status
			video.Status = status
			models.RDb.Del(context.Background(), define.VIDEO_STATUS+strconv.Itoa(lastStatus))
			models.RDb.Del(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(video.Vid))
			models.RDb.Del(context.Background(), define.DANMU_IDSET+strconv.Itoa(video.Vid))
			models.RDb.ZRem(context.Background(), define.USER_VIDEO_UPLOAD+strconv.Itoa(video.Uid), video.Vid)
			go func() {
				minIO.DelObject(videoUrl)
				oss.DelFile(coverUrl)
				result, err2 := models.RDb.ZRange(context.Background(), define.COMMENT_VIDEO+strconv.Itoa(video.Vid), 0, -1).Result()
				if err2 != nil {
					response.ResponseFailWithData(c, 500, "服务器错误", "")
					return
				}
				models.RDb.Del(context.Background(), define.COMMENT_VIDEO+strconv.Itoa(video.Vid))
				for _, v := range result {
					commentId, _ := strconv.Atoi(v)
					models.RDb.Del(context.Background(), define.COMMENT_REPLY+strconv.Itoa(commentId))
				}
			}()
			response.ResponseOK(c)
			return
		} else {
			response.ResponseFailWithData(c, 403, "无权限", "")
			return
		}
	}
}
func GetRandomVideos(c *gin.Context) {
	var count int64
	count = 11
	result, err := models.RDb.SRandMemberN(context.Background(), define.VIDEO_STATUS+strconv.Itoa(1), count).Result()
	if err != nil {
		response.ResponseFailWithData(c, 500, "服务器错误", "")
		return
	}
	ids := make([]int, len(result))
	for i, v := range result {
		id, _ := strconv.Atoi(v)
		ids[i] = id
	}
	videoList := getVideoByIds(ids, 1, int(count))
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ids), func(i, j int) {
		videoList[i], videoList[j] = videoList[j], videoList[i]
	})
	response.ResponseOKWithData(c, "", videoList)
	return
}

func CumulativeVideoForVisitor(c *gin.Context) {
	ids := c.Query("vids")
	idList := strings.Split(ids, ",")
	members, err := models.RDb.SMembers(context.Background(), define.VIDEO_STATUS+strconv.Itoa(1)).Result()
	var videoCumulative models.VideoCumulative
	if len(members) == 0 || err != nil {
		response.ResponseOKWithData(c, "", videoCumulative)
		return
	}
	var finalVidList []int
	for _, v := range members {
		isExist := helper.Contains(idList, v)
		if !isExist {
			finalId, _ := strconv.Atoi(v)
			finalVidList = append(finalVidList, finalId)
		}
	}
	minIndex := math.Min(float64(len(finalVidList)), float64(10))
	finalVidList = helper.GetShuffle(finalVidList)[:int(minIndex)]
	videoList := getVideoByIds(finalVidList, 1, 10)
	videoCumulative.Videos = videoList
	videoCumulative.Vids = finalVidList
	fmt.Println(len(members))
	fmt.Println(len(finalVidList))
	if len(members)-len(finalVidList)-10 > 0 {
		videoCumulative.More = true
	} else {
		videoCumulative.More = false
	}
	response.ResponseOKWithData(c, "", videoCumulative)
	return

}

func GetOneVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("vid"))
	videoGetVo := getVideoById(id)
	if videoGetVo == (models.VideoGetVo{}) {
		response.ResponseFailWithData(c, 404, "视频不见了", "")
		return
	}
	if videoGetVo.Video.Status != 1 {
		response.ResponseFailWithData(c, 404, "视频不见了", "")
		return
	}
	response.ResponseOKWithData(c, "", videoGetVo)
	return
}
func getVideoByIds(ids []int, pageNum int, pageSize int) []models.VideoGetVo {
	pageNum = (pageNum - 1) * pageSize
	if pageNum > len(ids) {
		return []models.VideoGetVo{}
	}
	endIndex := math.Min(float64(pageSize), float64(len(ids)))
	var videoList []models.VideoVo
	models.Db.Model(new(models.VideoVo)).
		Where("id in (?)", ids[pageNum:int(endIndex)]).
		Where("status != ?", 3).
		Find(&videoList)
	var mapList []models.VideoGetVo

	for _, v := range videoList {
		var videoMap models.VideoGetVo
		uid := v.Uid
		vid := v.Vid
		mcId := v.McId
		scId := v.ScId
		videoMap.Video = v
		user := getUserById(uid)
		videoStats := getVideoStatsById(vid)
		var category models.Category
		models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
		videoMap.User = user
		videoMap.Stats = videoStats
		videoMap.Category = category
		mapList = append(mapList, videoMap)
	}

	//sort.Slice(mapList, func(i, j int) bool {
	//	vidI := mapList[i].Video.Vid
	//	vidJ := mapList[j].Video.Vid
	//	return vidI < vidJ
	//})
	return mapList
}

func getVideoByIdList(ids []int) []models.VideoGetVo {
	var videoList []models.VideoVo
	models.Db.Model(new(models.VideoVo)).
		Where("id in (?)", ids).
		Where("status != ?", 3).
		Find(&videoList)
	var mapList []models.VideoGetVo

	for _, v := range videoList {
		var videoMap models.VideoGetVo
		uid := v.Uid
		vid := v.Vid
		mcId := v.McId
		scId := v.ScId
		videoMap.Video = v
		user := getUserById(uid)
		videoStats := getVideoStatsById(vid)
		var category models.Category
		models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
		videoMap.User = user
		videoMap.Stats = videoStats
		videoMap.Category = category
		mapList = append(mapList, videoMap)
	}

	//sort.Slice(mapList, func(i, j int) bool {
	//	vidI := mapList[i].Video.Vid
	//	vidJ := mapList[j].Video.Vid
	//	return vidI < vidJ
	//})
	return mapList
}
func getVideoById(id int) models.VideoGetVo {

	var video models.VideoVo

	result, err := models.RDb.Get(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(id)).Result()
	if err != nil {
		models.Db.Model(new(models.VideoVo)).
			Where("id = ?", id).
			Where("status != ?", 3).
			Find(&video)
		go func() {
			videoJson, _ := json.Marshal(video)
			models.RDb.Set(context.Background(), define.VIDEO_PREFIX+strconv.Itoa(id), videoJson, 0)
		}()
	} else {
		json.Unmarshal([]byte(result), &video)
	}

	var videoMap models.VideoGetVo
	uid := video.Uid
	vid := video.Vid
	mcId := video.McId
	scId := video.ScId
	videoMap.Video = video
	user := getUserById(uid)
	videoStats := getVideoStatsById(vid)
	var category models.Category
	models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
	videoMap.User = user
	videoMap.Stats = videoStats
	videoMap.Category = category

	return videoMap
}

func getVideosWithDataByIdsOrderbyDesc(ids []int, column string, pageNo int, pageSize int) []models.VideoGetVo {
	if column == "" {
		var videoList []models.VideoVo
		models.Db.Model(new(models.VideoVo)).
			Where("id in (?)", ids).
			Find(&videoList)
		var mapList []models.VideoGetVo

		for _, v := range videoList {
			var videoMap models.VideoGetVo
			if v.Status == 3 {
				videoMap.Video = v
				continue
			}
			uid := v.Uid
			vid := v.Vid
			mcId := v.McId
			scId := v.ScId
			videoMap.Video = v
			var user models.UserDto
			var videoStats models.VideoStats
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				user = getUserById(uid)
			}()
			go func() {
				defer wg.Done()
				videoStats = getVideoStatsById(vid)
			}()
			wg.Wait()
			var category models.Category
			models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
			videoMap.User = user
			videoMap.Stats = videoStats
			videoMap.Category = category
			mapList = append(mapList, videoMap)

		}
		return mapList
	} else if column == "upload_date" {
		var videoList []models.VideoVo
		models.Db.Model(new(models.VideoVo)).
			Where("id in (?)", ids).
			Order("upload_date desc").
			Offset((pageNo - 1) * pageSize).
			Limit(pageSize).
			Find(&videoList)
		var mapList []models.VideoGetVo
		for _, v := range videoList {
			var videoMap models.VideoGetVo
			if v.Status == 3 {
				videoMap.Video = v
				continue
			}
			uid := v.Uid
			vid := v.Vid
			mcId := v.McId
			scId := v.ScId
			videoMap.Video = v
			var user models.UserDto
			var videoStats models.VideoStats
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				user = getUserById(uid)
			}()
			go func() {
				defer wg.Done()
				videoStats = getVideoStatsById(vid)
			}()
			wg.Wait()
			var category models.Category
			models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
			videoMap.User = user
			videoMap.Stats = videoStats
			videoMap.Category = category
			mapList = append(mapList, videoMap)
		}
		return mapList
	} else {
		var videoStatsList []models.VideoStats
		models.Db.Model(new(models.VideoStats)).
			Where("vid in (?)", ids).
			Order(column + " desc").
			Offset((pageNo - 1) * pageSize).
			Limit(pageSize).
			Find(&videoStatsList)

		var mapList []models.VideoGetVo
		for _, v := range videoStatsList {
			var video models.VideoVo
			models.Db.Model(new(models.VideoVo)).
				Where("id = ?", v.Vid).
				Find(&video)

			var videoMap models.VideoGetVo
			if video.Status == 3 {
				videoMap.Video = video
				continue
			}
			uid := video.Uid
			vid := video.Vid
			mcId := video.McId
			scId := video.ScId
			videoMap.Video = video
			var user models.UserDto
			var videoStats models.VideoStats
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				user = getUserById(uid)
			}()
			go func() {
				defer wg.Done()
				videoStats = getVideoStatsById(vid)
			}()
			wg.Wait()
			var category models.Category
			models.Db.Model(new(models.Category)).Where("mc_id = ? and sc_id = ?", mcId, scId).Find(&category)
			videoMap.User = user
			videoMap.Stats = videoStats
			videoMap.Category = category
			mapList = append(mapList, videoMap)
		}
		return mapList
	}
}
