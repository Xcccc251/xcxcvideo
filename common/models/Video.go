package models

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Vid      int    `gorm:"column:id" json:"vid"`
	Uid      int    `gorm:"column:uid" json:"uid"`
	Title    string `gorm:"column:title" json:"title"`
	Type     int    `gorm:"column:type;default:1" json:"type"`
	Auth     int    `gorm:"column:auth;default:0" json:"auth"`
	Duration int    `gorm:"column:duration" json:"duration"`
	McId     string `gorm:"column:mc_id" json:"mcid"`
	ScId     string `gorm:"column:sc_id" json:"scid"`
	Tags     string `gorm:"column:tags" json:"tags"`
	Descr    string `gorm:"column:descr" json:"descr"`
	CoverUrl string `gorm:"column:cover_url" json:"coverUrl"`
	VideoUrl string `gorm:"column:video_url" json:"videoUrl"`
	Status   int    `gorm:"column:status;default:0" json:"status"`
}
type VideoVo struct {
	gorm.Model
	Vid        int    `gorm:"column:id" json:"vid"`
	Uid        int    `gorm:"column:uid" json:"uid"`
	Title      string `gorm:"column:title" json:"title"`
	Type       int    `gorm:"column:type;default:1" json:"type"`
	Auth       int    `gorm:"column:auth;default:0" json:"auth"`
	Duration   int    `gorm:"column:duration" json:"duration"`
	McId       string `gorm:"column:mc_id" json:"mcid"`
	ScId       string `gorm:"column:sc_id" json:"scid"`
	Tags       string `gorm:"column:tags" json:"tags"`
	Descr      string `gorm:"column:descr" json:"descr"`
	CoverUrl   string `gorm:"column:cover_url" json:"coverUrl"`
	VideoUrl   string `gorm:"column:video_url" json:"videoUrl"`
	Status     int    `gorm:"column:status;default:0" json:"status"`
	UploadDate MyTime `gorm:"column:created_at" json:"uploadDate"`
	DeleteDate MyTime `gorm:"column:deleted_at" json:"deleteDate"`
}
type VideoUploadInfoDto struct {
	gorm.Model
	Uid      int     `gorm:"column:id" json:"uid"`
	Hash     string  `gorm:"column:hash" json:"hash"`
	Title    string  `gorm:"column:title" json:"title"`
	Type     int     `gorm:"column:type;default:1" json:"type"`
	Auth     int     `gorm:"column:auth;default:0" json:"auth"`
	Duration float64 `gorm:"column:duration" json:"duration"`
	McId     string  `gorm:"column:mc_id" json:"mcid"`
	ScId     string  `gorm:"column:sc_id" json:"scid"`
	Tags     string  `gorm:"column:tags" json:"tags"`
	Descr    string  `gorm:"column:descr" json:"descr"`
	CoverUrl string  `gorm:"column:cover_url" json:"coverUrl"`
}
type VideoGetVo struct {
	Video    VideoVo    `json:"video"`
	User     UserDto    `json:"user"`
	Category Category   `json:"category"`
	Stats    VideoStats `json:"stats"`
}
type VideoCumulative struct {
	Videos []VideoGetVo `json:"videos"`
	Vids   []int        `json:"vids"`
	More   bool         `json:"more"`
}
type GetUserWorksDto struct {
	Count int          `json:"count"`
	List  []VideoGetVo `json:"list"`
}

func (Video) TableName() string {
	return "video"
}
func (VideoVo) TableName() string {
	return "video"
}
