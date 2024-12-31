package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id          int     `gorm:"column:id" json:"uid"`
	Username    string  `gorm:"column:username" json:"username"`
	Password    string  `gorm:"column:password" json:"password"`
	Nickname    string  `gorm:"column:nickname" json:"nickname"`
	Avatar      string  `gorm:"column:avatar" json:"avatar"`
	BackGround  string  `gorm:"column:background" json:"background"`
	Gender      int     `gorm:"column:gender;default:2" json:"gender"`
	Description string  `gorm:"column:description" json:"description"`
	Exp         int     `gorm:"column:exp;default:0" json:"exp"`
	Coin        float64 `gorm:"column:coin;default:0" json:"coin"`
	Vip         int     `gorm:"column:vip;default:0" json:"vip"`
	State       int     `gorm:"column:state;default:0" json:"state"`
	Role        int     `gorm:"column:role;default:0" json:"role"`
	Auth        int     `gorm:"column:auth;default:0" json:"auth"`
	AuthMsg     string  `gorm:"column:auth_msg" json:"authMsg"`
}
type UserVo struct {
	gorm.Model
	Id          int     `gorm:"column:id" json:"uid"`
	Username    string  `gorm:"column:username" json:"username"`
	Password    string  `gorm:"column:password" json:"password"`
	Nickname    string  `gorm:"column:nickname" json:"nickname"`
	Avatar      string  `gorm:"column:avatar" json:"avatar"`
	BackGround  string  `gorm:"column:background" json:"background"`
	Gender      int     `gorm:"column:gender;default:2" json:"gender"`
	Description string  `gorm:"column:description" json:"description"`
	Exp         int     `gorm:"column:exp;default:0" json:"exp"`
	Coin        float64 `gorm:"column:coin;default:0" json:"coin"`
	Vip         int     `gorm:"column:vip;default:0" json:"vip"`
	State       int     `gorm:"column:state;default:0" json:"state"`
	Role        int     `gorm:"column:role;default:0" json:"role"`
	Auth        int     `gorm:"column:auth;default:0" json:"auth"`
	AuthMsg     string  `gorm:"column:auth_msg" json:"authMsg"`
	CreateDate  string  `gorm:"column:created_at" json:"createDate"`
	UpdateDate  string  `gorm:"column:updated_at" json:"updateDate"`
}
type UserDto struct {
	Id           int     `gorm:"column:id" json:"uid"`
	Nickname     string  `gorm:"column:nickname" json:"nickname"`
	Avatar       string  `gorm:"column:avatar" json:"avatar_url"`
	BackGround   string  `gorm:"column:background" json:"bg_url"`
	Gender       int     `gorm:"column:gender;default:2" json:"gender"`
	Description  string  `gorm:"column:description" json:"description"`
	Exp          int     `gorm:"column:exp;default:0" json:"exp"`
	Coin         float64 `gorm:"column:coin;default:0" json:"coin"`
	Vip          int     `gorm:"column:vip;default:0" json:"vip"`
	State        int     `gorm:"column:state;default:0" json:"state"`
	Auth         int     `gorm:"column:auth;default:0" json:"auth"`
	AuthMsg      string  `gorm:"column:auth_msg" json:"authMsg"`
	VideoCount   int     `gorm:"column:video_count" json:"video_count"`
	FollowsCount int     `gorm:"column:follows_count" json:"follows_count"`
	FansCount    int     `gorm:"column:fans_count" json:"fans_count"`
	LoveCount    int     `gorm:"column:love_count" json:"love_count"`
	PlayCount    int     `gorm:"column:play_count" json:"play_count"`
}
type UserLoginOrRegisterDto struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmedPassword"`
}
type LoginRsp struct {
	User  UserDto `json:"user"`
	Token string  `json:"token"`
}

func (User) TableName() string {
	return "user"
}
func (UserVo) TableName() string {
	return "user"
}
