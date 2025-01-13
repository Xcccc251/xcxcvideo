package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetDanmuList(c *gin.Context) {
	vid, _ := strconv.Atoi(c.Param("vid"))
	idSet, _ := models.RDb.SMembers(context.Background(), define.DANMU_IDSET+strconv.Itoa(vid)).Result()
	var idSetInt []int
	for _, id := range idSet {
		idInt, _ := strconv.Atoi(id)
		idSetInt = append(idSetInt, idInt)
	}
	var danmuList []models.DanmuVo
	models.Db.Model(new(models.DanmuVo)).Where("id in ?", idSetInt).Where("state=?", 1).Find(&danmuList)
	response.ResponseOKWithData(c, "获取弹幕成功", danmuList)
	return
}

func DelDanmu(c *gin.Context) {
	userId, _ := c.Get("userId")
	isAdmin, _ := c.Get("isAdmin")
	id, _ := strconv.Atoi(c.PostForm("id"))
	db := models.Db.Model(new(models.Danmu)).Where("id=?", id).Where("state != ?", 3)
	var count int64
	db.Count(&count)
	if count == 0 {
		response.ResponseFailWithData(c, 0, "弹幕不存在", nil)
	}
	var danmu models.Danmu
	db.Find(&danmu)
	var video models.Video
	models.Db.Model(new(models.Video)).Where("id=?", danmu.Vid).Find(&video)
	if isAdmin.(bool) || danmu.Uid == userId.(int) || video.Uid == userId.(int) {
		db.Update("state", 3)
		UpdateVideoStats(danmu.Vid, "danmu", false, 1)
		models.RDb.SRem(context.Background(), define.DANMU_IDSET+strconv.Itoa(danmu.Vid), danmu.Id)
		response.ResponseOKWithData(c, "删除弹幕成功", nil)
		return
	}
	response.ResponseFailWithData(c, 0, "无权限", nil)
	return
}
