package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Id         int    `gorm:"column:id" json:"id"`
	UserId     int    `gorm:"column:user_id" json:"userId"`
	AnotherId  int    `gorm:"column:another_id" json:"anotherId"`
	IsDeleted  int    `gorm:"column:is_deleted;default:0" json:"isDeleted"`
	Unread     int    `gorm:"column:unread;default:0" json:"unread"`
	LatestTime MyTime `gorm:"column:latest_time" json:"latestTime"`
}
type ChatGetVo struct {
	Chat       Chat        `json:"chat"`
	User       UserDto     `json:"user"`
	ChatDetail interface{} `json:"chatDetail"`
}
type ChatDetail struct {
	Id         int    `gorm:"column:id" json:"id"`
	UserId     int    `gorm:"column:user_id" json:"userId"`
	AnotherId  int    `gorm:"column:another_id" json:"anotherId"`
	Content    string `gorm:"column:content" json:"content"`
	UserDel    int    `gorm:"column:user_del" json:"userDel"`
	AnotherDel int    `gorm:"column:another_del" json:"anotherDel"`
	Withdraw   int    `gorm:"column:withdraw" json:"withdraw"`
	Time       MyTime `gorm:"column:time" json:"time"`
}
type RecentListGetVo struct {
	List []ChatGetVo `json:"list"`
	More bool        `json:"more"`
}
type ChatDetailGetVo struct {
	List []ChatDetail `json:"list"`
	More bool         `json:"more"`
}

func (Chat) TableName() string {
	return "chat"
}

func (ChatDetail) TableName() string {
	return "chat_detailed"
}
