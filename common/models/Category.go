package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	McId   string `gorm:"column:mc_id" json:"mcId"`
	ScId   string `gorm:"column:sc_id" json:"scId"`
	McName string `gorm:"column:mc_name" json:"mcName"`
	ScName string `gorm:"column:sc_name" json:"scName"`
	Descr  string `gorm:"column:descr" json:"descr"`
	RcmTag string `gorm:"column:rcm_tag" json:"rcmTag"`
}
type CategoryDto struct {
	gorm.Model
	McId   string             `gorm:"column:mc_id" json:"mcId"`
	McName string             `gorm:"column:mc_name" json:"mcName"`
	ScList []ChildrenCategory `gorm:"column:sc_list" json:"scList"`
}
type ChildrenCategory struct {
	McId   string   `gorm:"column:mc_id" json:"mcId"`
	ScId   string   `gorm:"column:sc_id" json:"scId"`
	McName string   `gorm:"column:mc_name" json:"mcName"`
	ScName string   `gorm:"column:sc_name" json:"scName"`
	Descr  string   `gorm:"column:descr" json:"descr"`
	RacTag []string `gorm:"column:rac_tag" json:"racTag"`
}

func (m *Category) TableName() string {
	return "category"
}
