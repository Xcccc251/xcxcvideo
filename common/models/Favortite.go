package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	Fid         int    `gorm:"column:fid" json:"fid"`
	Uid         int    `gorm:"column:uid" json:"uid"`
	Type        int    `gorm:"column:type" json:"type"`
	Visible     int    `gorm:"column:visible" json:"visible"`
	Cover       string `gorm:"column:cover" json:"cover"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	Count       int    `gorm:"column:count" json:"count"`
	IsDelete    int    `gorm:"column:is_delete" json:"is_delete"`
}

type FavoriteVideo struct {
	gorm.Model
	Id       int    `gorm:"column:id" json:"id"`
	Vid      int    `gorm:"column:vid" json:"vid"`
	Fid      int    `gorm:"column:fid" json:"fid"`
	Time     MyTime `gorm:"column:time" json:"time"`
	IsRemove int    `gorm:"column:is_remove" json:"is_remove"`
}

func (Favorite) TableName() string {
	return "favorite"
}

func (FavoriteVideo) TableName() string {
	return "favorite_video"
}
