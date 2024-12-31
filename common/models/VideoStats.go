package models

import "gorm.io/gorm"

type VideoStats struct {
	gorm.Model
	Vid     int `gorm:"column:vid" json:"vid"`
	Play    int `gorm:"column:play" json:"play"`
	Danmu   int `gorm:"column:danmu" json:"danmu"`
	Good    int `gorm:"column:good" json:"good"`
	Bad     int `gorm:"column:bad" json:"bad"`
	Coin    int `gorm:"column:coin" json:"coin"`
	Collect int `gorm:"column:collect" json:"collect"`
	Share   int `gorm:"column:share" json:"share"`
	Comment int `gorm:"column:comment" json:"comment"`
}

func (VideoStats) TableName() string {
	return "video_stats"
}
