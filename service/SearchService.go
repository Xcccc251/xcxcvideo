package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/es"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SearchHotList(c *gin.Context) {
	result, _ := models.RDb.ZRevRange(context.Background(), define.SEARCH_WORD, 0, 9).Result()
	var hotSearchList []models.HotSearch
	for _, v := range result {
		hotSearch := models.HotSearch{}
		hotSearch.Content = v
		score, _ := models.RDb.ZScore(context.Background(), define.SEARCH_WORD, v).Result()
		hotSearch.Score = score
		lastScore := getScoreByKeyword(v)
		if lastScore == 0 {
			hotSearch.Type = 1
			if hotSearch.Score > 3 {
				hotSearch.Type = 2
			}
		} else if score-lastScore > 3 {
			hotSearch.Type = 2
		}
		hotSearchList = append(hotSearchList, hotSearch)
	}
	sort.Slice(hotSearchList, func(i, j int) bool {
		return hotSearchList[i].Score > hotSearchList[j].Score
	})
	response.ResponseOKWithData(c, "", hotSearchList)
	return
}

func getScoreByKeyword(keyword string) float64 {
	result, err := models.RDb.Get(context.Background(), define.HOT_SEARCH_WORDS).Result()
	if err != nil {
		return 0
	}
	var hotSearchWordsList []models.HotSearchWords
	json.Unmarshal([]byte(result), &hotSearchWordsList)
	for _, v := range hotSearchWordsList {
		if v.Keyword == keyword {
			return v.Score
		}
	}
	return 0
}
func AddSearchWord(c *gin.Context) {
	keyword := c.PostForm("keyword")
	formattedString := formatString(keyword)
	if len(formattedString) == 0 {
		response.ResponseOKWithData(c, "", nil)
		return
	}
	if len(formattedString) < 2 || len(formattedString) > 30 {
		response.ResponseOKWithData(c, "", formattedString)
		return
	}
	go func() {
		// 检查有序集合中是否包含该值
		ctx := context.Background()
		_, err := models.RDb.ZScore(ctx, define.SEARCH_WORD, formattedString).Result()
		if err == redis.Nil {
			fmt.Println("值不存在，执行 ZIncrBy")
			// 值不存在，执行 ZIncrBy
			models.RDb.ZIncrBy(ctx, define.SEARCH_WORD, 1, formattedString)
			es.AddSearchWord(formattedString)
		} else if err != nil {
			// 处理其他错误
			response.ResponseOKWithData(c, "内部错误", nil)
			return
		} else {
			// 值存在，可以执行其他操作，例如增加分数
			models.RDb.ZIncrBy(ctx, define.SEARCH_WORD, 1, formattedString)

		}
	}()
	response.ResponseOKWithData(c, "", formattedString)
	return
}
func formatString(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\x{4e00}-\x{9fff}\x{3041}-\x{309f}\x{30a1}-\x{30ff}]+`)
	return strings.Trim(re.ReplaceAllString(input, ""), "")
}

func GetSearchWord(c *gin.Context) {
	keyword := c.Query("keyword")
	finalKeyword, err := url.QueryUnescape(keyword)
	if err != nil {
		response.ResponseOKWithData(c, "内部错误", nil)
		return
	}
	response.ResponseOKWithData(c, "", es.GetSearchWord(finalKeyword))
	return
}

func GetSearchCount(c *gin.Context) {
	keyword := c.Query("keyword")
	finalKeyword, err := url.QueryUnescape(keyword)
	if err != nil {
		response.ResponseOKWithData(c, "内部错误", nil)
		return
	}
	sw := sync.WaitGroup{}
	sw.Add(1)
	videoCount := 0
	var userCount int64

	go func() {
		videoCount = es.GetSearchVideoCount(finalKeyword)
		sw.Done()
	}()
	go func() {
		models.Db.Model(new(models.User)).Where("nickname like ?", "%"+finalKeyword+"%").Count(&userCount)
	}()
	sw.Wait()
	var result []int
	result = append(result, videoCount)
	result = append(result, int(userCount))
	response.ResponseOKWithData(c, "", result)
	return

}

func GetSearchVideo(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(define.DEFAULT_PAGE_NUM)))
	finalKeyword, err := url.QueryUnescape(keyword)
	if err != nil {
		response.ResponseOKWithData(c, "内部错误", nil)
		return
	}
	vids := es.GetSearchVideo(finalKeyword, page, 30, true)
	list := getVideoByIdList(vids)
	response.ResponseOKWithData(c, "", list)
	return
}

func GetSearchUser(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(define.DEFAULT_PAGE_NUM)))
	finalKeyword, err := url.QueryUnescape(keyword)
	if err != nil {
		response.ResponseOKWithData(c, "内部错误", nil)
		return
	}
	var users []models.User
	models.Db.Model(new(models.User)).Where("nickname like ?", "%"+finalKeyword+"%").Offset((page - 1) * 30).Limit(30).Find(&users)
	var userDtos []models.UserDto
	copier.Copy(&userDtos, &users)
	response.ResponseOKWithData(c, "", userDtos)
	return
}
