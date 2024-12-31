package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Vid        int    `gorm:"column:id" json:"vid"`
	Uid        int    `gorm:"column:uid" json:"uid"`
	Title      string `gorm:"column:title" json:"title"`
	Type       int    `gorm:"column:type;default:1" json:"type"`
	Auth       int    `gorm:"column:auth;default:0" json:"auth"`
	Duration   int    `gorm:"column:duration" json:"duration"`
	McId       int    `gorm:"column:mc_id" json:"mcId"`
	ScId       int    `gorm:"column:sc_id" json:"scId"`
	Tags       string `gorm:"column:tags" json:"tags"`
	Descr      string `gorm:"column:descr" json:"descr"`
	CoverUrl   string `gorm:"column:cover_url" json:"coverUrl"`
	VideoUrl   string `gorm:"column:video_url" json:"videoUrl"`
	Status     int    `gorm:"column:status;default:0" json:"status"`
	UploadDate string `gorm:"column:upload_date" json:"uploadDate"`
	DeleteDate string `gorm:"column:deleted_at" json:"deleteDate"`
}

func (Video) TableName() string {
	return "video"
}
