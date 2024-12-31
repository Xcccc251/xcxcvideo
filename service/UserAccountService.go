package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/helper"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var USER_ERROR_CODE = 403

func Registser(c *gin.Context) {
	var userRegisterDto models.UserRegisterDto
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
	newIdStr := strconv.Itoa(define.ID_PREFIX + int(lastId) + 1)
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
