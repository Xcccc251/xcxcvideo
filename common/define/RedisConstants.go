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
	USER_VIDEO_UPLOAD     = "xcxc:user:video:upload:"
	VIDEOSTATS_PREFIX     = "xcxc:videoStats:"
	CATEGORYLIST          = "xcxc:categoryList"
)
