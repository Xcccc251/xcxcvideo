package define

import "time"

var (
	USER_PREFIX           = "xcxc:user:"
	TOKEN_PREFIX          = "xcxc:token:"
	TOKEN_TTL             = time.Hour * 24
	MSG_UNREAD            = "xcxc:msg:unread:"
	MSG_UNREAD_TTL        = time.Hour * 24
	FAVORITE_PREFIX       = "xcxc:favorite:"
	FAVORITE_VIDEO_PREFIX = "xcxc:favorite:video:"
	DEFAULT_TTL           = time.Hour * 2
	USER_LIKE_COMMENT     = "xcxc:user:like:comment:"
	USER_DISLIKE_COMMENT  = "xcxc:user:dislike:comment:"
	USER_VIDEO_UPLOAD     = "xcxc:user:video:upload:"
	VIDEOSTATS_PREFIX     = "xcxc:videoStats:"
)
