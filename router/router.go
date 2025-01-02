package router

import (
	"XcxcVideo/common/middlewares"
	"XcxcVideo/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	user := r.Group("/user")
	{
		user.POST("/account/register", service.Registser)
		user.POST("/account/login", service.Login)
		user.GET("/personal/info", middlewares.AuthUserCheck(), service.GetUserInfo)
		user.GET("/info/get-one", service.GetUserInfoById)
		user.POST("/info/update", middlewares.AuthUserCheck(), service.UpdateUserInfo)
		user.POST("/avatar/update", middlewares.AuthUserCheck(), service.UpdateAvatar)
	}
	admin := r.Group("/admin")
	{
		admin.POST("/account/login", service.AdminLogin)
		admin.GET("/personal/info", middlewares.AuthAdminCheck(), service.GetAdminInfo)
	}
	msg_unread := r.Group("/msg-unread")
	{
		msg_unread.GET("/all", middlewares.AuthUserCheck(), service.GetMsgUnread)
	}
	favorite := r.Group("/favorite")
	{
		favorite.GET("/get-all/user", middlewares.AuthUserCheck(), service.GetAllFavoritesForUser)
	}
	comment := r.Group("/comment")
	{
		comment.GET("/get-like-and-dislike", middlewares.AuthUserCheck(), service.GetUserLikeAndDislike)
	}
	category := r.Group("/category")
	{
		category.GET("/getall", service.GetCategoryList)
	}

	return r
}
