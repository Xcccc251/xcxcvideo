package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Id       int    `gorm:"column:id" json:"id"`
	Vid      int    `gorm:"column:vid" json:"vid"`
	Uid      int    `gorm:"column:uid" json:"uid"`
	RootId   int    `gorm:"column:root_id" json:"rootId"`
	ParentId int    `gorm:"column:parent_id" json:"parentId"`
	ToUserId int    `gorm:"column:to_user_id" json:"toUserId"`
	Content  string `gorm:"column:content" json:"content"`
	Love     int    `gorm:"column:love" json:"love"`
	Bad      int    `gorm:"column:bad" json:"bad"`
	IsTop    int    `gorm:"column:is_top" json:"isTop"`
}
type CommentVo struct {
	gorm.Model
	Id         int    `gorm:"column:id" json:"id"`
	Vid        int    `gorm:"column:vid" json:"vid"`
	Uid        int    `gorm:"column:uid" json:"uid"`
	RootId     int    `gorm:"column:root_id" json:"rootId"`
	ParentId   int    `gorm:"column:parent_id" json:"parentId"`
	ToUserId   int    `gorm:"column:to_user_id" json:"toUserId"`
	Content    string `gorm:"column:content" json:"content"`
	Love       int    `gorm:"column:love" json:"love"`
	Bad        int    `gorm:"column:bad" json:"bad"`
	IsTop      int    `gorm:"column:is_top" json:"isTop"`
	CreateTime MyTime `gorm:"column:created_at" json:"createTime"`
}

func (Comment) TableName() string {
	return "comment"
}
func (CommentVo) TableName() string {
	return "comment"
}
