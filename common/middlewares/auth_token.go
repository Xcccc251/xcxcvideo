package middlewares

import (
	"SecKill/common/define"
	"SecKill/common/models"
	"context"
	"github.com/gin-gonic/gin"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.Next()
		} else {
			// 刷新过期时间
			models.RDb.Expire(context.Background(), define.LOGIN_USER_KEY+token, define.LOGIN_USER_TTL)
			c.Next()
		}

	}
}
