package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	websocketServer "XcxcVideo/websocket"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
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

		if comment.ToUserId == comment.Uid {
			models.RDb.ZAdd(context.Background(), define.REPLY_ZSET+strconv.Itoa(comment.ToUserId), &redis.Z{
				Member: comment.Id,
				Score:  float64(time.Now().Unix()),
			})

			msgMap := make(map[string]interface{})
			msgMap["type"] = "接收"

			websocketServer.SendMessage(comment.ToUserId, "reply", msgMap)
			addOneUnreadMsg(comment.ToUserId, "reply")
		}

	}()
	return

}

func GetCommentTreeByVid(c *gin.Context) {
	var commentGetVo models.CommentGetVo
	vid, _ := strconv.Atoi(c.Query("vid"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	typeOf, _ := strconv.Atoi(c.Query("type"))
	result, _ := models.RDb.ZCard(context.Background(), define.COMMENT_VIDEO+strconv.Itoa(vid)).Result()
	if offset > int(result) {
		commentGetVo.More = false
		commentGetVo.Comments = make([]models.CommentTree, 0)
		response.ResponseOKWithData(c, "", commentGetVo)
		return

	}
	if offset+10 > int(result) {
		commentGetVo.More = false
	} else {
		commentGetVo.More = true
	}

	commentVoList := getRootCommentsByVid(vid, offset, typeOf)
	commentTreeList := make([]models.CommentTree, len(commentVoList))
	for i, v := range commentVoList {
		commentTreeList[i] = buildCommentTree(v, 0, 2)
	}
	commentGetVo.Comments = commentTreeList
	response.ResponseOKWithData(c, "", commentGetVo)
	return

}
func GetRootComments(c *gin.Context) {
	vid, _ := strconv.Atoi(c.Query("vid"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	typeOfQuery, _ := strconv.Atoi(c.Query("type"))
	commentVoList := getRootCommentsByVid(vid, offset, typeOfQuery)
	response.ResponseOKWithData(c, "", commentVoList)
	return
}
func getRootCommentsByVid(vid int, offset int, typeOfQuery int) []models.CommentVo {
	var rootIdsSet []string
	var commentVoList []models.CommentVo
	if typeOfQuery == 1 {
		fmt.Println("热度查询")
		models.Db.Model(new(models.CommentVo)).
			Where("root_id = ?", 0).
			Where("vid = ?", vid).
			Order("(love-bad) desc").
			Offset(offset).Limit(10).
			Find(&commentVoList)
	} else {
		fmt.Println("时间查询")
		rootIdsSet, _ = models.RDb.ZRange(context.Background(), define.COMMENT_VIDEO+strconv.Itoa(vid), int64(offset), int64(offset+9)).Result()
		models.Db.Model(new(models.CommentVo)).Where("id in ?", rootIdsSet).Order("id desc").Find(&commentVoList)
	}
	return commentVoList
}
func getChildCommentsByRootId(rootId int, vid int, start int, stop int) []models.CommentVo {
	result, err := models.RDb.ZRange(context.Background(), define.COMMENT_REPLY+strconv.Itoa(rootId), int64(start), int64(stop)).Result()
	if err != nil {
		return make([]models.CommentVo, 0)
	}
	var commentVoList []models.CommentVo
	models.Db.Model(new(models.CommentVo)).Where("id in ?", result).Find(&commentVoList)
	return commentVoList

}
func buildCommentTree(comment models.CommentVo, start int, stop int) models.CommentTree {
	var commentTree models.CommentTree
	copier.Copy(&commentTree, &comment)
	commentTree.User = getUserById(comment.Uid)
	commentTree.ToUser = getUserById(comment.ToUserId)
	if comment.RootId == 0 {
		result, err := models.RDb.ZCard(context.Background(), define.COMMENT_REPLY+strconv.Itoa(comment.Id)).Result()
		if err != nil {
			return models.CommentTree{}
		}
		commentTree.Count = int(result)

		childComments := getChildCommentsByRootId(comment.Id, comment.Vid, start, stop)
		childCommentTrees := make([]models.CommentTree, len(childComments))
		for i, v := range childComments {

			childCommentTrees[i] = buildCommentTree(v, start, stop)
		}
		commentTree.Replies = childCommentTrees

	}
	return commentTree

}

func updateLikeOrDisLike(id int, addLike bool) {
	if addLike {
		var comment models.Comment
		db := models.Db.Model(new(models.Comment)).Where("id = ?", id)
		db.Find(&comment)
		db.Update("love", gorm.Expr("love+1"))
		if comment.Bad > 0 {
			db.Update("bad", gorm.Expr("bad-1"))
		}
	} else {
		var comment models.Comment
		db := models.Db.Model(new(models.Comment)).Where("id = ?", id)
		db.Find(&comment)
		db.Update("bad", gorm.Expr("bad+1"))
		if comment.Love > 0 {
			db.Update("love", gorm.Expr("love-1"))
		}
	}

}

func updateCommentColumn(id int, column string, icr bool, count int) {
	db := models.Db.Model(new(models.Comment)).Where("id = ?", id)
	if icr {
		db.Update(column, gorm.Expr(column+"+"+strconv.Itoa(count)))
	} else {
		var preCount int64
		db.Select(column).Find(&preCount)
		if preCount > int64(count) {
			db.Update(column, gorm.Expr(column+"-"+strconv.Itoa(count)))
		} else {
			db.Update(column, 0)
		}
	}
}
