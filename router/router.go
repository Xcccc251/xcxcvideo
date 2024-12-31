package router

import (
	"XcxcVideo/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	user := r.Group("/user")
	{
		user.POST("/account/register", service.Registser)
	}

	return r
}
