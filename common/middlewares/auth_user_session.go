package middlewares

import (
	"SecKill/common/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthUserCheckBySession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user") == nil {
			c.Abort()
			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
			return
		}
		user := session.Get("user")
		c.Set("userId", user.(*models.User).Id)
		c.Set("user", user)
		c.Next()

	}
}
