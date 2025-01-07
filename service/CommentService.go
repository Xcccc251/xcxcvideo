package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func AddComment(c *gin.Context) {
	userId, _ := c.Get("userId")
	vid, _ := strconv.Atoi(c.PostForm("vid"))
	rootId, _ := strconv.Atoi(c.PostForm("root_id"))
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	toUserId, _ := strconv.Atoi(c.PostForm("to_user_id"))
	content := c.PostForm("content")
	if content == "" || len(content) == 0 {
		response.ResponseFailWithData(c, 400, "评论内容不能为空", "")
		return
	}
	if len(content) > 2000 {
		response.ResponseFailWithData(c, 400, "评论内容过长", "")
		return
	}
	var comment models.Comment
	comment.Vid = vid
	comment.RootId = rootId
	comment.ParentId = parentId
	comment.ToUserId = toUserId
	comment.Uid = userId.(int)
	comment.Content = content
	models.Db.Model(new(models.Comment)).Create(&comment)
	updateVideoStats(vid, "comment", true, 1)
	response.ResponseOKWithData(c, "评论成功", nil)
	go func() {
		if rootId != 0 {
			models.RDb.ZAdd(context.Background(), define.COMMENT_REPLY+strconv.Itoa(rootId), &redis.Z{
				Member: comment.Id,
				Score:  float64(time.Now().Unix()),
			})
		} else {
			models.RDb.ZAdd(context.Background(), define.COMMENT_VIDEO+strconv.Itoa(vid), &redis.Z{
				Member: comment.Id,
				Score:  float64(time.Now().Unix()),
			})
		}

		if comment.ToUserId != comment.Uid {
			models.RDb.ZAdd(context.Background(), define.REPLY_ZSET+strconv.Itoa(comment.ToUserId), &redis.Z{
				Member: comment.Id,
				Score:  float64(time.Now().Unix()),
			})
		}
		//todo 未读消息
	}()
	return

}
