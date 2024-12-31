package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

func GetUserLikeAndDislike(c *gin.Context) {
	userId, _ := c.Get("userId")
	var likeList []string
	var dislikeList []string
	sw := sync.WaitGroup{}
	sw.Add(2)
	go func() {
		likeList = redisUtil.GetSet(define.USER_LIKE_COMMENT + strconv.Itoa(userId.(int)))
		sw.Done()
	}()

	go func() {
		dislikeList = redisUtil.GetSet(define.USER_DISLIKE_COMMENT + strconv.Itoa(userId.(int)))
		sw.Done()
	}()
	sw.Wait()

	response.ResponseOKWithData(c, "获取成功", gin.H{
		"userLike":    likeList,
		"userDislike": dislikeList,
	})

}
