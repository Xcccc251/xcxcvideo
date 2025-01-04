package service

import (
	"XcxcVideo/common/define"
	"XcxcVideo/common/models"
	"XcxcVideo/common/response"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"strings"
)

func GetCategoryList(c *gin.Context) {
	var categoryList []models.CategoryDto
	result, err := models.RDb.Get(context.Background(), define.CATEGORYLIST).Result()
	if err == nil {
		json.Unmarshal([]byte(result), &categoryList)
		response.ResponseOKWithData(c, "获取成功", categoryList)
		return
	}
	var dbCategoryList []models.Category
	models.Db.Model(new(models.Category)).Find(&dbCategoryList)
	var categoryDtoMap = map[string]models.CategoryDto{}
	for _, v := range dbCategoryList {
		rcmTagString := v.RcmTag
		rcmTagList := strings.Split(rcmTagString, "\n")
		if _, exists := categoryDtoMap[v.McId]; exists == false {
			var categoryDto models.CategoryDto
			categoryDto.McId = v.McId
			categoryDto.McName = v.McName
			categoryDto.ScList = []models.ChildrenCategory{}
			categoryDtoMap[v.McId] = categoryDto
		}
		var childrenCategory models.ChildrenCategory
		copier.Copy(&childrenCategory, &v)
		childrenCategory.RcmTag = rcmTagList
		categoryDtoTemp := categoryDtoMap[v.McId]
		categoryDtoTemp.ScList = append(categoryDtoTemp.ScList, childrenCategory)
		categoryDtoMap[v.McId] = categoryDtoTemp

	}
	for _, v := range define.DEFAULT_CATEGORY_ORDER {
		if _, exists := categoryDtoMap[v]; exists == true {
			categoryList = append(categoryList, categoryDtoMap[v])
		}
	}
	go func() {
		models.RDb.Set(context.Background(), define.CATEGORYLIST, categoryList, 0)
	}()
	response.ResponseOKWithData(c, "获取成功", categoryList)
	return

}
