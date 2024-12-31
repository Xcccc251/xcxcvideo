package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	Fid         int    `gorm:"column:id" json:"fid"`
	Uid         int    `gorm:"column:uid" json:"uid"`
	Type        int    `gorm:"column:type" json:"type"`
	Visible     int    `gorm:"column:visible" json:"visible"`
	Cover       string `gorm:"column:cover" json:"cover"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	Count       int    `gorm:"column:count" json:"count"`
	IsDelete    int    `gorm:"column:is_delete" json:"is_delete"`
}

func (Favorite) TableName() string {
	return "favorite"
}
