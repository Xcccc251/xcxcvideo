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
	USER_LIKE_COMMENT     = "xcxc:user_like_comment:"
	USER_DISLIKE_COMMENT  = "xcxc:user_dislike_comment:"
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
	DANMU_IDSET           = "xcxc:danmu_idset:"
	COMMENT_VIDEO         = "xcxc:comment_video:"
	COMMENT_REPLY         = "xcxc:comment_reply:"
	REPLY_ZSET            = "xcxc:reply_zset:"
	USER_VIDEO_HISTORY    = "xcxc:user_video_history:"
	SEARCH_WORD           = "xcxc:search_word:"
	HOT_SEARCH_WORDS      = "xcxc:hot_search_words:"
	CHAT_ZSET             = "xcxc:chat_zset:"
	CHAT_DETAILED_ZSET    = "xcxc:chat_detailed_zset:"
	WHISPER_KEY           = "xcxc:whisper:"
	LOVE_VIDEO            = "xcxc:love_video:"
	BELOVED_VIDEO_SET     = "xcxc:beloved_video_set:"

	USER_FAVORITES      = "xcxc:user_favorites:"
	USER_FAVORITE_VIDEO = "xcxc:user_favorite_video:"
)
