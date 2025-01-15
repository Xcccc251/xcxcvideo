package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"XcxcVideo/common/oss"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"path"
	"strconv"
	"strings"
)

func GetUserInfoById(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("uid"))
	userDto := getUserById(userId)
	response.ResponseOKWithData(c, "", userDto)
	return
}
func UpdateUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	nickname := c.PostForm("nickname")
	desc := c.PostForm("description")
	gender, _ := strconv.Atoi(c.PostForm("gender"))
	nickname = strings.Trim(nickname, " ")
	if nickname == "" {
		response.ResponseFailWithData(c, 500, "昵称不能为空", nil)
		return
	}
	if len(nickname) > 30 || len(desc) > 100 {
		response.ResponseFailWithData(c, 500, "输入字符过长", nil)
		return
	}
	if nickname == "账号已注销" {
		response.ResponseFailWithData(c, 500, "非法昵称", nil)
		return
	}
	var count int64
	models.Db.Model(new(models.User)).
		Where("nickname = ?", nickname).
		Where("id != ?", userId.(int)).Count(&count)
	if count > 0 {
		response.ResponseFailWithData(c, 500, "昵称已被使用", nil)
		return
	}
	models.Db.Model(new(models.User)).
		Where("id = ?", userId.(int)).
		Update("nickname", nickname).
		Update("description", desc).
		Update("gender", gender)
	//todo 更新es
	go func() {
		models.RDb.Del(context.Background(), define.USER_PREFIX+strconv.Itoa(userId.(int)))
		var user models.UserVo
		models.Db.Model(new(models.UserVo)).Where("id = ?", userId.(int)).First(&user)
		redisUtil.Set(define.USER_PREFIX+strconv.Itoa(userId.(int)), user, define.DEFAULT_TTL)
	}()

	response.ResponseOK(c)
	return
}

func UpdateAvatar(c *gin.Context) {
	userId, _ := c.Get("userId")
	file, _ := c.FormFile("file")
	filename := file.Filename
	ext := path.Ext(filename)
	objectName := helper.GetUUID() + ext
	filePath, err := file.Open()
	if err != nil {
		response.ResponseFailWithData(c, 500, "上传失败", nil)
		return
	}
	uploadFilePath, err := oss.UploadFile(objectName, filePath)
	if err != nil {
		response.ResponseFailWithData(c, 500, "上传失败", nil)
		return
	}
	models.Db.Model(new(models.User)).Where("id = ?", userId).Update("avatar", uploadFilePath)
	go func() {
		models.RDb.Del(context.Background(), define.USER_PREFIX+strconv.Itoa(userId.(int)))
		var user models.UserVo
		models.Db.Model(new(models.UserVo)).Where("id = ?", userId.(int)).First(&user)
		redisUtil.Set(define.USER_PREFIX+strconv.Itoa(userId.(int)), user, define.DEFAULT_TTL)
	}()
	response.ResponseOKWithData(c, "", uploadFilePath)
	return
}
func getUserById(userId int) models.UserDto {
	var userDto models.UserDto
	var user models.UserVo
	userResult, err := models.RDb.Get(context.Background(), define.USER_PREFIX+strconv.Itoa(userId)).Result()
	json.Unmarshal([]byte(userResult), &user)
	if err != nil {
		models.Db.Model(new(models.UserVo)).Where("id = ?", userId).First(&user)
		redisUtil.Set(define.USER_PREFIX+strconv.Itoa(userId), user, define.DEFAULT_TTL)
		return userDto
	}

	copier.Copy(&userDto, &user)
	if user.State == 2 {
		userDto.Nickname = "账号已注销"
		userDto.Avatar = define.DEFAULT_AVATAR_URL
		userDto.BackGround = define.DEFAULT_BACKGROUND_URL
		userDto.Description = "账号已注销"
		userDto.Gender = define.GENDER_UNKOWN
		return userDto
	}

	videoList := redisUtil.GetSet(define.USER_VIDEO_UPLOAD + strconv.Itoa(userId))
	if len(videoList) == 0 {
		return userDto
	}
	var videoCount int
	var loveCount int
	var playCount int
	videoCount = len(videoList)
	loveCount, playCount = processVideoStats(videoList)
	userDto.VideoCount = videoCount
	userDto.LoveCount = loveCount
	userDto.PlayCount = playCount
	return userDto
}

func IsAdmin(userId int) bool {
	var user models.User
	models.Db.Model(new(models.User)).Where("id = ?", userId).First(&user)
	return user.Role != 0
}

func GetUserWorks(c *gin.Context) {
	var getUserWorks models.GetUserWorksDto
	userId, _ := strconv.Atoi(c.Query("uid"))
	rule, _ := strconv.Atoi(c.Query("rule"))
	pageNo, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("quantity"))

	result, err := models.RDb.SMembers(context.Background(), define.USER_VIDEO_UPLOAD+strconv.Itoa(userId)).Result()
	if err != nil {
		getUserWorks.List = make([]models.VideoGetVo, 0)
		getUserWorks.Count = 0
		response.ResponseOKWithData(c, "", getUserWorks)
		return
	}
	var ids []int
	for _, idStr := range result {
		id, _ := strconv.Atoi(idStr)
		ids = append(ids, id)
	}
	getUserWorks.Count = len(ids)
	switch rule {
	case 1:
		getUserWorks.List = getVideosWithDataByIdsOrderbyDesc(ids, "upload_date", pageNo, pageSize)
		break
	case 2:
		getUserWorks.List = getVideosWithDataByIdsOrderbyDesc(ids, "play", pageNo, pageSize)
		break
	case 3:
		getUserWorks.List = getVideosWithDataByIdsOrderbyDesc(ids, "good", pageNo, pageSize)
		break
	default:
		getUserWorks.List = getVideosWithDataByIdsOrderbyDesc(ids, "upload_date", pageNo, pageSize)

	}
	response.ResponseOKWithData(c, "", getUserWorks)
	return

}
