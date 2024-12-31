package models

import "gorm.io/gorm"

type MsgUnread struct {
	gorm.Model
	UId     int `gorm:"column:id" json:"uid"`
	Reply   int `gorm:"column:reply;default:0" json:"reply"`
	At      int `gorm:"column:at;default:0" json:"at"`
	Love    int `gorm:"column:love;default:0" json:"love"`
	System  int `gorm:"column:system;default:0" json:"system"`
	Whisper int `gorm:"column:whisper;default:0" json:"whisper"`
	Dynamic int `gorm:"column:dynamic;default:0" json:"dynamic"`
}

func (MsgUnread) TableName() string {
	return "msg_unread"
}
