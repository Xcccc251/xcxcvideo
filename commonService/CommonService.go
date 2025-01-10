package commonService

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/redisUtil"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

func UpdateChat(from int, to int) bool {
	//查询对方聊天是否在线
	key := define.WHISPER_KEY + strconv.Itoa(to) + ":" + strconv.Itoa(from)
	result, _ := models.RDb.Exists(context.Background(), key).Result()
	sw := sync.WaitGroup{}
	sw.Add(2)
	go func() {
		var chat1 models.Chat
		db := models.Db.Model(new(models.Chat)).Where("user_id = ? and another_id = ?", to, from)
		db.Update("latest_time", models.MyTime(time.Now()))
		db.Update("is_deleted", 0)
		db.Find(&chat1)
		models.RDb.ZAdd(context.Background(), define.CHAT_ZSET+strconv.Itoa(from), &redis.Z{
			Member: chat1.Id,
			Score:  float64(time.Now().Unix()),
		})
		sw.Done()
	}()

	go func() {
		var chat2 models.Chat
		var isExist int64
		db := models.Db.Model(new(models.Chat)).Where("user_id = ? and another_id = ?", from, to)
		db.Count(&isExist)
		if result != 0 {
			if isExist == 0 {
				db.Create(&models.Chat{
					UserId:     from,
					AnotherId:  to,
					LatestTime: models.MyTime(time.Now()),
				})
			} else {
				db.Update("latest_time", models.MyTime(time.Now()))
				db.Update("is_deleted", 0)
				db.Find(&chat2)
			}
			models.RDb.ZAdd(context.Background(), define.CHAT_ZSET+strconv.Itoa(to), &redis.Z{
				Member: chat2.Id,
				Score:  float64(time.Now().Unix()),
			})
			sw.Done()
		} else {
			if isExist == 0 {
				db.Create(&models.Chat{
					UserId:     from,
					AnotherId:  to,
					Unread:     1,
					LatestTime: models.MyTime(time.Now()),
				})
			} else {
				db.Update("latest_time", models.MyTime(time.Now()))
				db.Update("is_deleted", 0)
				db.Update("unread", gorm.Expr("unread + ?", 1))
			}
			addOneUnreadMsg(to, "whisper")
			models.RDb.ZAdd(context.Background(), define.CHAT_ZSET+strconv.Itoa(to), &redis.Z{
				Member: chat2.Id,
				Score:  float64(time.Now().Unix()),
			})
			sw.Done()
		}

	}()
	sw.Wait()
	return result != 0

}

func addOneUnreadMsg(userId int, column string) {
	models.Db.Model(new(models.MsgUnread)).Where("uid = ?", userId).Update(column, gorm.Expr(column+"+1"))
	models.RDb.Del(context.Background(), define.MSG_UNREAD+strconv.Itoa(userId))
}

func GetUserById(userId int) models.UserDto {
	var userDto models.UserDto
	var user models.UserVo
	userResult, err := models.RDb.Get(context.Background(), define.USER_PREFIX+strconv.Itoa(userId)).Result()
	json.Unmarshal([]byte(userResult), &user)
	if err != nil {
		models.Db.Model(new(models.UserVo)).Where("id = ?", userId).First(&user)
		redisUtil.Set(define.USER_PREFIX+strconv.Itoa(userId), user, define.DEFAULT_TTL)
		return userDto
	}

	copier.Copy(&userDto, &user)
	if user.State == 2 {
		userDto.Nickname = "账号已注销"
		userDto.Avatar = define.DEFAULT_AVATAR_URL
		userDto.BackGround = define.DEFAULT_BACKGROUND_URL
		userDto.Description = "账号已注销"
		userDto.Gender = define.GENDER_UNKOWN
		return userDto
	}

	videoList := redisUtil.GetSet(define.USER_VIDEO_UPLOAD + strconv.Itoa(userId))
	if len(videoList) == 0 {
		return userDto
	}
	var videoCount int
	var loveCount int
	var playCount int
	videoCount = len(videoList)
	loveCount, playCount = processVideoStats(videoList)
	userDto.VideoCount = videoCount
	userDto.LoveCount = loveCount
	userDto.PlayCount = playCount
	return userDto
}
func processVideoStats(videoList []string) (loveCount int, playCount int) {
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, 10) // 控制最大并发数为10
	)

	for _, videoId := range videoList {
		wg.Add(1)
		sem <- struct{}{}
		vid, _ := strconv.Atoi(videoId)
		//创建新线程查询视频统计数据
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			videoStats := getVideoStatsById(vid)
			loveCount += videoStats.Good
			playCount += videoStats.Play
		}()
	}
	wg.Wait()
	return loveCount, playCount

}
func getVideoStatsById(vid int) models.VideoStats {
	var videoStats models.VideoStats

	videoId := strconv.Itoa(vid)
	result, err := models.RDb.Get(context.Background(), define.VIDEOSTATS_PREFIX+videoId).Result()
	if err == nil {
		json.Unmarshal([]byte(result), &videoStats)
		models.RDb.Expire(context.Background(), define.VIDEOSTATS_PREFIX+videoId, define.DEFAULT_TTL)
		return videoStats
	}

	models.Db.Model(new(models.VideoStats)).Where("vid = ?", videoId).Find(&videoStats)
	go func() {
		videoStatsJson, _ := json.Marshal(videoStats)
		models.RDb.Set(context.Background(), define.VIDEOSTATS_PREFIX+videoId, videoStatsJson, define.DEFAULT_TTL)
	}()
	return videoStats

}
