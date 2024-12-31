package middlewares

import (
	"SecKill/common/define"
	"SecKill/common/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func AuthUserCheckByRedis() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.Abort()
			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
			return
		}
		result, err := models.RDb.Get(context.Background(), define.LOGIN_USER_KEY+token).Result()
		if err != nil {
			c.Abort()
			c.JSON(401, gin.H{"code": 401, "msg": "请先登录"})
			return
		}
		// 刷新过期时间
		models.RDb.Expire(context.Background(), define.LOGIN_USER_KEY+token, define.LOGIN_USER_TTL)
		var user models.User
		json.Unmarshal([]byte(result), &user)
		var userDto models.UserDto
		copier.Copy(&userDto, &user)
		c.Set("userId", user.Id)
		c.Set("user", userDto)
		c.Next()
	}
}
