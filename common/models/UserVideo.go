package models

import "gorm.io/gorm"

type UserVideo struct {
	gorm.Model
	Id       int    `gorm:"column:id" json:"id"`
	Uid      int    `gorm:"column:uid" json:"uid"`
	Vid      int    `gorm:"column:vid" json:"vid"`
	Play     int    `gorm:"column:play" json:"play"`
	Love     int    `gorm:"column:love" json:"love"`
	Unlove   int    `gorm:"column:unlove" json:"unlove"`
	Coin     int    `gorm:"column:coin" json:"coin"`
	Collect  int    `gorm:"column:collect" json:"collect"`
	PlayTime MyTime `gorm:"column:play_time" json:"playTime"`
	LoveTime MyTime `gorm:"column:love_time" json:"loveTime"`
	CoinTime MyTime `gorm:"column:coin_time" json:"coinTime"`
}

func (UserVideo) TableName() string {
	return "user_video"
}
