package helper

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"strconv"
	"time"
)

func Set(key string, value interface{}, duration time.Duration) {
	valueJson, _ := json.Marshal(value)
	models.RDb.Set(context.Background(), key, valueJson, duration)
}
func SetWithLogicalExpire(key string, value interface{}, duration time.Duration) {
	cacheValue := struct {
		Value      interface{}
		ExpireTime int64
	}{
		Value:      value,
		ExpireTime: time.Now().Add(duration).Unix(),
	}
	cacheValueJson, _ := json.Marshal(cacheValue)
	models.RDb.Set(context.Background(), key, cacheValueJson, 0)
}

type QueryFunc[T any] func(id int, value *T) error

func QueryWithPassThrough[T any](keyPrefix string, id int, value T, duration time.Duration) error {
	idStr := strconv.Itoa(id)
	key := keyPrefix + idStr
	jsonStr, err := models.RDb.Get(context.Background(), key).Result()

	if err == redis.Nil {
		if err2 := models.Db.Model(new(T)).Where("id = ?", id).First(&value).Error; err2 != nil {
			models.RDb.Set(context.Background(), key, "", define.CACHE_NULL_TTL)
			return errors.New("不存在的值")
		}
		jsonStr2, _ := json.Marshal(value)
		models.RDb.Set(context.Background(), key, jsonStr2, duration)
		return nil
	}
	if jsonStr == "" {
		return errors.New("不存在的值")
	}
	json.Unmarshal([]byte(jsonStr), &value)
	return nil

}

func QueryWithLogicalExpire[T any](cacheKeyPrefix string, lockKeyPrefix string, id int, value T, duration time.Duration) error {
	cacheKey := cacheKeyPrefix + strconv.Itoa(id)
	var cacheT struct {
		Value      T
		ExpireTime int64
	}
	cacheJson, err := models.RDb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		lockKey := lockKeyPrefix + strconv.Itoa(id)
		if tryLock(lockKey) {
			defer UnLock(lockKey)
			if models.Db.Model(new(T)).Where("id = ?", id).First(&value).Error != nil {
				models.RDb.Set(context.Background(), cacheKey, "", define.CACHE_NULL_TTL)
				return errors.New("不存在的值")
			}
			cacheT.Value = value
			cacheT.ExpireTime = time.Now().Add(duration).Unix()
			cacheTJson, _ := json.Marshal(cacheT)
			models.RDb.Set(context.Background(), cacheKey, cacheTJson, 0)
			return nil
		} else {
			return errors.New("系统繁忙，请稍后重试！")
		}
	} else if err != nil {
		return errors.New("系统繁忙，请稍后重试！")
	}
	if cacheJson == "" {
		models.RDb.Set(context.Background(), cacheKey, "", define.CACHE_NULL_TTL)
		return errors.New("不存在的值！")
	}
	json.Unmarshal([]byte(cacheJson), &cacheT)
	currentTime := time.Now().Unix()
	if currentTime < cacheT.ExpireTime {
		copier.Copy(&value, &cacheT.Value)
		return nil
	}
	lockKey := define.LOCK_SHOP_KEY + strconv.Itoa(id)
	if tryLock(lockKey) {
		go func() {
			defer UnLock(lockKey)
			if models.Db.Model(new(T)).Where("id = ?", id).First(&value).Error != nil {
				models.RDb.Set(context.Background(), cacheKey, "", define.CACHE_NULL_TTL)
			}
			cacheT.Value = value
			cacheT.ExpireTime = time.Now().Add(duration).Unix()
			cacheTJson, _ := json.Marshal(cacheT)
			models.RDb.Set(context.Background(), cacheKey, cacheTJson, 0)
		}()
	}
	copier.Copy(&value, &cacheT.Value)
	return nil
}

func tryLock(key string) bool {
	return models.RDb.SetNX(context.Background(), key, 1, define.LOCK_SHOP_TTL).Val()
}
func UnLock(key string) {
	models.RDb.Del(context.Background(), key)
}
