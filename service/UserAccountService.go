package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var USER_ERROR_CODE = 403

func Registser(c *gin.Context) {
	var userRegisterDto models.UserLoginOrRegisterDto
	if err := c.ShouldBindJSON(&userRegisterDto); err != nil {
		response.ResponseFailWithData(c, http.StatusInternalServerError, "参数错误", nil)
		return
	}
	userRegisterDto.Username = strings.ReplaceAll(userRegisterDto.Username, " ", "")
	if userRegisterDto.Username == "" {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名不能为空", nil)
		return
	}
	if userRegisterDto.Password == "" || userRegisterDto.ConfirmedPassword == "" {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "密码不能为空", nil)
		return
	}
	if len(userRegisterDto.Username) > 50 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名长度不能超过50", nil)
		return
	}
	if len(userRegisterDto.Password) > 50 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "密码长度不能超过50", nil)
		return
	}
	if userRegisterDto.Password != userRegisterDto.ConfirmedPassword {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "密码不一致", nil)
		return
	}
	var count int64
	db := models.Db.Model(new(models.User)).
		Where("username = ?", userRegisterDto.Username).
		Where("state != ?", 2)
	db.Count(&count)
	if count > 0 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户已存在", nil)
		return
	}
	var lastId int64
	models.Db.Model(new(models.User)).Select("id").Order("id desc").Limit(1).Find(&lastId)
	var user models.User
	newId := define.ID_PREFIX + int(lastId) + 1
	newIdStr := define.NICKNAME_PREFIX + strconv.Itoa(define.ID_PREFIX+int(lastId)+1)
	user.Username = userRegisterDto.Username
	user.Password, _ = helper.GetBcryptPassword(userRegisterDto.Password)
	user.Description = "这个人很懒，什么都没有留下"
	user.Avatar = define.DEFAULT_AVATAR_URL
	user.BackGround = define.DEFAULT_BACKGROUND_URL
	user.Nickname = newIdStr
	user.Id = newId
	if err := models.Db.Create(&user).Error; err != nil {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "注册失败", nil)
		return
	}
	//msgUnreadMapper.insert(new MsgUnread(new_user.getUid(),0,0,0,0,0,0));
	//favoriteMapper.insert(new Favorite(null, new_user.getUid(), 1, 1, null, "默认收藏夹", "", 0, null));
	//esUtil.addUser(new_user);
	//TODO
	response.ResponseOKWithData(c, "注册成功,欢迎加入我们", nil)
	return

}

func Login(c *gin.Context) {
	var userLoginDto models.UserLoginOrRegisterDto
	if err := c.ShouldBindJSON(&userLoginDto); err != nil {
		response.ResponseFailWithData(c, http.StatusInternalServerError, "参数错误", nil)
		return
	}
	var dbUser models.UserVo
	var count int64
	db := models.Db.Model(new(models.UserVo)).Where("username = ?", userLoginDto.Username)
	if db.Count(&count); count == 0 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名或密码错误", nil)
		return
	}
	db.First(&dbUser)
	if !helper.AnalysisBcryptPassword(dbUser.Password, userLoginDto.Password) {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名或密码错误", nil)
		return
	}
	if dbUser.State == 1 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "账号封禁中", nil)
		return
	}
	token, _ := helper.GenerateToken(dbUser.Id)
	userJson, _ := json.Marshal(dbUser)
	idStr := strconv.Itoa(dbUser.Id)
	err := models.RDb.SetEX(context.Background(), define.TOKEN_PREFIX+idStr, token, define.TOKEN_TTL).Err()
	if err != nil {
		fmt.Println(err)
	}
	models.RDb.SetEX(context.Background(), define.USER_PREFIX+idStr, userJson, define.TOKEN_TTL)
	var userDto models.UserDto
	copier.Copy(&userDto, &dbUser)
	var loginRsp models.LoginRsp
	loginRsp.Token = token
	loginRsp.User = userDto
	response.ResponseOKWithData(c, "欢迎回来", loginRsp)
	return

}

func GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	var user models.UserVo
	models.Db.Model(new(models.UserVo)).Where("id = ?", userId).First(&user)
	if user.State == 2 {
		response.ResponseFailWithData(c, 404, "账号已注销", nil)
		return
	}
	if user.State == 1 {
		response.ResponseFailWithData(c, 403, "账号封禁中", nil)
		return
	}

	var userDto models.UserDto
	copier.Copy(&userDto, &user)

	videoList := redisUtil.GetSet(define.USER_VIDEO_UPLOAD + strconv.Itoa(userId.(int)))
	if len(videoList) == 0 {
		response.ResponseOKWithData(c, "", userDto)
		return
	}
	var videoCount int
	var loveCount int
	var playCount int
	videoCount = len(videoList)
	loveCount, playCount = processVideoStats(videoList)
	userDto.VideoCount = videoCount
	userDto.LoveCount = loveCount
	userDto.PlayCount = playCount
	response.ResponseOKWithData(c, "", userDto)
	return

}

func AdminLogin(c *gin.Context) {
	var adminLoginDto models.AdminLoginDto
	if err := c.ShouldBindJSON(&adminLoginDto); err != nil {
		response.ResponseFailWithData(c, http.StatusInternalServerError, "参数错误", nil)
		return
	}
	var dbAdmin models.UserVo
	err := models.Db.Model(new(models.UserVo)).Where("username = ?", adminLoginDto.Username).Find(&dbAdmin).Error
	if err != nil {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名或密码错误", nil)
		return
	}
	if !helper.AnalysisBcryptPassword(dbAdmin.Password, adminLoginDto.Password) {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "用户名或密码错误", nil)
		return
	}
	if dbAdmin.Role == 0 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "您不是管理员，无权访问", nil)
		return
	}
	if dbAdmin.State == 1 {
		response.ResponseFailWithData(c, USER_ERROR_CODE, "账号封禁中", nil)
		return
	}
	token, _ := helper.GenerateToken(dbAdmin.Id)
	dbAdminJson, _ := json.Marshal(dbAdmin)
	models.RDb.Set(context.Background(), define.TOKEN_PREFIX+strconv.Itoa(dbAdmin.Id), token, define.TOKEN_TTL)
	models.RDb.Set(context.Background(), define.ADMIN_PREFIX+strconv.Itoa(dbAdmin.Id), dbAdminJson, define.ADMIN_LOGIN_TTL)
	var userDto models.UserDto
	copier.Copy(&userDto, &dbAdmin)
	var loginRsp models.LoginRsp
	loginRsp.Token = token
	loginRsp.User = userDto
	response.ResponseOKWithData(c, "欢迎回来", loginRsp)
	return
}
func GetAdminInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	result, err := models.RDb.Get(context.Background(), define.ADMIN_PREFIX+strconv.Itoa(userId.(int))).Result()
	var admin models.UserVo
	if err != nil {
		models.Db.Model(new(models.UserVo)).Where("id = ?", userId).First(&admin)
		adminJson, _ := json.Marshal(admin)
		redisUtil.Set(define.ADMIN_PREFIX+strconv.Itoa(userId.(int)), adminJson, define.ADMIN_LOGIN_TTL)
	} else {
		json.Unmarshal([]byte(result), &admin)
	}
	var userDto models.UserDto
	copier.Copy(&userDto, &admin)
	response.ResponseOKWithData(c, "", userDto)
}

// 优化遍历
func processVideoStats(videoList []string) (loveCount int, playCount int) {
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, 10) // 控制最大并发数为10
	)

	for _, videoId := range videoList {
		wg.Add(1)
		sem <- struct{}{}
		vid, _ := strconv.Atoi(videoId)
		//创建新线程查询视频统计数据
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			videoStats := getVideoStatsById(vid)
			loveCount += videoStats.Good
			playCount += videoStats.Play
		}()
	}
	wg.Wait()
	return loveCount, playCount

}
