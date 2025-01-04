package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
)

func GetTotalVideoCount(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("vstatus"))
	result, _ := models.RDb.SMembers(context.Background(), define.VIDEO_STATUS+strconv.Itoa(status)).Result()
	count := len(result)
	response.ResponseOKWithData(c, "获取成功", count)
	return

}
func GetOneReviewVideo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("vid"))
	response.ResponseOKWithData(c, "获取成功", getVideoById(id))
	return
}
func GetReviewVideo(c *gin.Context) {
	status, _ := strconv.Atoi(c.Query("vstatus"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(define.DEFAULT_PAGE_NUM)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("quantity", strconv.Itoa(define.DEFAULT_PAGE_SIZE)))
	//count := models.RDb.SCard(context.Background(), define.VIDEO_STATUS+strconv.Itoa(status)).Val()
	members := models.RDb.SMembers(context.Background(), define.VIDEO_STATUS+strconv.Itoa(status)).Val()
	var ids []int
	for _, member := range members {
		id, _ := strconv.Atoi(member)
		ids = append(ids, id)
	}
	var data []models.VideoGetVo
	if len(members) != 0 {
		data = getVideoByIds(ids, page, pageSize)
	}
	sort.Slice(data, func(i, j int) bool {
		vidI := data[i].Video.Vid
		vidJ := data[j].Video.Vid
		return vidI < vidJ
	})
	response.ResponseOKWithData(c, "获取成功", data)
	return
}
