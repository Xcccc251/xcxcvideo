package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/redisUtil"
	"XcxcVideo/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

func GetUserLikeAndDislike(c *gin.Context) {
	userId, _ := c.Get("userId")
	likeList := []int{}
	dislikeList := []int{}
	var likeResult []string
	var dislikeResult []string
	sw := sync.WaitGroup{}
	sw.Add(2)
	go func() {
		likeResult = redisUtil.GetSet(define.USER_LIKE_COMMENT + strconv.Itoa(userId.(int)))
		for _, v := range likeResult {
			id, _ := strconv.Atoi(v)
			likeList = append(likeList, id)
		}
		sw.Done()
	}()

	go func() {
		dislikeResult = redisUtil.GetSet(define.USER_DISLIKE_COMMENT + strconv.Itoa(userId.(int)))
		for _, v := range dislikeResult {
			id, _ := strconv.Atoi(v)
			dislikeList = append(dislikeList, id)
		}
		sw.Done()
	}()
	sw.Wait()

	response.ResponseOKWithData(c, "获取成功", gin.H{
		"userLike":    likeList,
		"userDislike": dislikeList,
	})

}

func LikeOrDisLikeComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	isLikeStr := c.PostForm("isLike")
	isLike, _ := strconv.ParseBool(isLikeStr)
	isSetStr := c.PostForm("isSet")
	isSet, _ := strconv.ParseBool(isSetStr)
	userId, _ := c.Get("userId")
	likeExist, _ := models.RDb.SIsMember(context.Background(), define.USER_LIKE_COMMENT+strconv.Itoa(userId.(int)), id).Result()
	dislikeExist, _ := models.RDb.SIsMember(context.Background(), define.USER_DISLIKE_COMMENT+strconv.Itoa(userId.(int)), id).Result()
	if isLike && isSet {
		if likeExist {
			response.ResponseOK(c)
			return
		}
		if dislikeExist {
			go func() {
				models.RDb.SRem(context.Background(), define.USER_DISLIKE_COMMENT+strconv.Itoa(userId.(int)), id)
			}()
			go func() {
				updateLikeOrDisLike(id, true)
			}()
		} else {
			go func() {
				updateCommentColumn(id, "love", true, 1)
			}()
		}
		models.RDb.SAdd(context.Background(), define.USER_LIKE_COMMENT+strconv.Itoa(userId.(int)), id)
	} else if isLike {
		if !likeExist {
			response.ResponseOK(c)
			return
		}
		go func() {
			models.RDb.SRem(context.Background(), define.USER_LIKE_COMMENT+strconv.Itoa(userId.(int)), id)
		}()
		go func() {
			updateCommentColumn(id, "love", false, 1)
		}()
	} else if isSet {
		if dislikeExist {
			response.ResponseOK(c)
			return
		}
		if likeExist {
			go func() {
				models.RDb.SRem(context.Background(), define.USER_LIKE_COMMENT+strconv.Itoa(userId.(int)), id)
			}()
			go func() {
				updateLikeOrDisLike(id, false)
			}()
		} else {
			go func() {
				updateCommentColumn(id, "bad", true, 1)
			}()
		}
		models.RDb.SAdd(context.Background(), define.USER_DISLIKE_COMMENT+strconv.Itoa(userId.(int)), id)
	} else {
		if !dislikeExist {
			response.ResponseOK(c)
			return
		}
		go func() {
			models.RDb.SRem(context.Background(), define.USER_DISLIKE_COMMENT+strconv.Itoa(userId.(int)), id)
		}()
		go func() {
			updateCommentColumn(id, "bad", false, 1)
		}()
	}
	response.ResponseOK(c)
	return

}
