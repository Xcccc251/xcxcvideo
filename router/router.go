package router

import (
	"XcxcVideo/common/middlewares"
	"XcxcVideo/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CorsMiddleWare())
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
		msg_unread.POST("/clear", middlewares.AuthUserCheck(), service.ClearUnreadMsg)
	}
	favorite := r.Group("/favorite")
	{
		favorite.GET("/get-all/user", middlewares.AuthUserCheck(), service.GetAllFavoritesForUser)
	}
	comment := r.Group("/comment")
	{
		comment.GET("/get-like-and-dislike", middlewares.AuthUserCheck(), service.GetUserLikeAndDislike)
		comment.POST("/add", middlewares.AuthUserCheck(), service.AddComment)
		comment.GET("/getroot", service.GetRootComments)
		comment.GET(("/get"), service.GetCommentTreeByVid)
		comment.POST("love-or-not", middlewares.AuthUserCheck(), service.LikeOrDisLikeComment)
	}
	category := r.Group("/category")
	{
		category.GET("/getall", service.GetCategoryList)
	}
	video := r.Group("/video")
	{
		video.GET("ask-chunk", middlewares.AuthUserCheck(), service.AskCurrentChunkByHash)
		video.POST("upload-chunk", middlewares.AuthUserCheck(), service.UploadVideoChunk)
		video.GET("cancel-upload", middlewares.AuthUserCheck(), service.CancelUpload)
		video.POST("add", middlewares.AuthUserCheck(), service.UploadVideo)
		video.GET("/getone", service.GetVideoById)
		video.POST("change/status", middlewares.AuthAdminCheck(), service.ChangeVideoStatus)
		video.GET("/random/visitor", service.GetRandomVideos)
		video.POST("play/user", middlewares.AuthUserCheck(), service.UserPlayVideo)
	}
	review_video := r.Group("/review/video")
	{
		review_video.GET("/total", middlewares.AuthAdminCheck(), service.GetTotalVideoCount)
		review_video.GET("/getpage", middlewares.AuthAdminCheck(), service.GetReviewVideo)
		review_video.GET("/getone", middlewares.AuthAdminCheck(), service.GetOneReviewVideo)
	}

	search := r.Group("/search")
	{
		//search.GET("/hot/get", service.SearchHotList)
		search.POST("/word/add", service.AddSearchWord)
		search.GET("/word/get", service.GetSearchWord)
		search.GET("/hot/get", service.SearchHotList)
		search.GET("/count", service.GetSearchCount)
		search.GET("/video/only-pass", service.GetSearchVideo)
	}

	chat := r.Group("/msg/chat")
	{
		chat.GET("/recent-list", middlewares.AuthUserCheck(), service.GetRecentLIst)
	}

	return r
}
