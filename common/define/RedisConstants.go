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
)
