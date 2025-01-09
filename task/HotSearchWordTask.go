package task

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sort"
)

func UpdateHotSearchWord() {
	//todo 热搜词条算法
	fmt.Println("每小时更新热搜词")
	result, err := models.RDb.ZRange(context.Background(), define.SEARCH_WORD, 0, -1).Result()
	if len(result) == 0 || err != nil {
		return
	}
	count := len(result)
	var total float64
	for _, v := range result {
		score, _ := models.RDb.ZScore(context.Background(), define.SEARCH_WORD, v).Result()
		total += score
	}
	var hotSearchWordList []models.HotSearchWords
	for i := 0; i < count; i++ {
		score, _ := models.RDb.ZScore(context.Background(), define.SEARCH_WORD, result[i]).Result()
		//归一化分数
		models.RDb.ZAdd(context.Background(), define.SEARCH_WORD, &redis.Z{Score: (score / total) * float64(count), Member: result[i]})
		hotSearchWordList = append(hotSearchWordList, models.HotSearchWords{Keyword: result[i], Score: (score / total) * float64(count)})
	}

	if err != nil {
		fmt.Println("json marshal error:", err)
		return
	}
	sort.Slice(hotSearchWordList, func(i, j int) bool {
		return hotSearchWordList[i].Score > hotSearchWordList[j].Score
	})
	hotSearchWordListJson, err := json.Marshal(hotSearchWordList)

	models.RDb.Set(context.Background(), define.HOT_SEARCH_WORDS, string(hotSearchWordListJson), 0)
	s, err := models.RDb.Get(context.Background(), define.HOT_SEARCH_WORDS).Result()
	var hotSearchWords []models.HotSearchWords
	json.Unmarshal([]byte(s), &hotSearchWords)
	fmt.Println(hotSearchWords[0].Score)

}
