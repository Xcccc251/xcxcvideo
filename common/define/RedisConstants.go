package define

import "time"

var (
	USER_PREFIX           = "xcxc:user:"
	ADMIN_PREFIX          = "xcxc:admin:"
	TOKEN_PREFIX          = "xcxc:token:"
	TOKEN_TTL             = time.Hour * 24
	ADMIN_LOGIN_TTL       = time.Hour * 2
	MSG_UNREAD            = "xcxc:msg:unread:"
	MSG_UNREAD_TTL        = time.Hour * 24
	FAVORITE_PREFIX       = "xcxc:favorite:"
	FAVORITE_VIDEO_PREFIX = "xcxc:favorite:video:"
	DEFAULT_TTL           = time.Hour * 2
	USER_LIKE_COMMENT     = "xcxc:user:like:comment:"
	USER_DISLIKE_COMMENT  = "xcxc:user:dislike:comment:"
	VIDEOSTATS_PREFIX     = "xcxc:videoStats:"
	CATEGORYLIST          = "xcxc:categoryList"
	VIDEO_PREFIX          = "xcxc:video:"
	VIDEO_STATUS          = "xcxc:video_status:"
	VIDEO_STATUS_0        = "xcxc:video_status:0"
	VIDEO_STATUS_1        = "xcxc:video_status:1"
	VIDEO_STATUS_2        = "xcxc:video_status:2"
	VIDEO_STATUS_3        = "xcxc:video_status:3"
	VIDEOSTATS            = "xcxc:videoStats:"
	USER_VIDEO_UPLOAD     = "xcxc:user_video_upload:"
)
