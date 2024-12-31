package define

import "time"

var (
	LOGIN_CODE_KEY     = "seckill:login:code:"
	LOGIN_CODE_TTL     = 60 * 5 * time.Second
	LOGIN_USER_KEY     = "seckill:login:token:"
	LOGIN_USER_TTL     = 60 * 60 * time.Second
	CACHE_NULL_TTL     = 2 * 60 * time.Second
	CACHE_SHOP_TTL     = 30 * 60 * time.Second
	CACHE_SHOPTYPE_TTL = 30 * 60 * time.Second
	CACHE_SHOP_KEY     = "seckill:cache:shop:"
	CACHE_SHOPTYPE_KEY = "seckill:cache:shoptype:"
	LOCK_SHOP_KEY      = "seckill:lock:shop:"
	LOCK_SHOP_TTL      = 10 * time.Second
	SECKILL_STOCK_KEY  = "seckill:stock:"
	BLOG_LIKED_KEY     = "seckill:blog:liked:"
	FEED_KEY           = "seckill:feed:"
	SHOP_GEO_KEY       = "seckill:shop:geo:"
	USER_SIGN_KEY      = "seckill:sign:"
	ICR_KEY            = "seckill:icr:"
	LOCK_PREFIX        = "seckill:lock:"
)
