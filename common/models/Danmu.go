package models

import "gorm.io/gorm"

type Danmu struct {
	gorm.Model
	Id         int     `gorm:"column:id" json:"id"`
	Vid        int     `gorm:"column:vid" json:"vid"`
	Uid        int     `gorm:"column:uid" json:"uid"`
	Content    string  `gorm:"column:content" json:"content"`
	Fontsize   int     `gorm:"column:fontsize" json:"fontsize"`
	Mode       int     `gorm:"column:mode" json:"mode"`
	Color      string  `gorm:"column:color" json:"color"`
	TimePoint  float64 `gorm:"column:time_point" json:"timePoint"`
	State      int     `gorm:"column:state" json:"state"`
	CreateDate MyTime  `gorm:"column:created_at" json:"createDate"`
}

type DanmuVo struct {
	gorm.Model
	Id         int     `gorm:"column:id" json:"id"`
	Vid        int     `gorm:"column:vid" json:"vid"`
	Uid        int     `gorm:"column:uid" json:"uid"`
	Content    string  `gorm:"column:content" json:"content"`
	Fontsize   int     `gorm:"column:fontsize" json:"fontsize"`
	Mode       int     `gorm:"column:mode" json:"mode"`
	Color      string  `gorm:"column:color" json:"color"`
	TimePoint  float64 `gorm:"column:time_point" json:"timePoint"`
	State      int     `gorm:"column:state" json:"state"`
	CreateDate MyTime  `gorm:"column:created_at" json:"createDate"`
}

func (Danmu) TableName() string {
	return "danmu"
}
func (DanmuVo) TableName() string {
	return "danmu"
}

func (d *Danmu) BeforeCreate(tx *gorm.DB) (err error) {
	d.State = 1
	return
}
