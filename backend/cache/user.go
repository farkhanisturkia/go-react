package cache

import (
	"strconv"
	"time"

	"go-react-vue/backend/database"
	"go-react-vue/backend/models"
	"go-react-vue/backend/pkg/redis"
)

const (
	UserListCachePrefix  = "users:list:page:"
	UserTotalCacheKey    = "users:total"
	UserListKeysSet      = "users:list:keys"
	UserCachePrefix      = "user:id:"
	CacheTTL             = 5 * time.Minute
	UserTTL              = 10 * time.Minute
)

func GetUsersTotalCount() int64 {
	val, err := redis.Client.Get(redis.Ctx, UserTotalCacheKey).Result()
	if err == nil {
		if total, err := strconv.ParseInt(val, 10, 64); err == nil {
			return total
		}
	}

	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	redis.Client.Set(redis.Ctx, UserTotalCacheKey, strconv.FormatInt(total, 10), CacheTTL)

	return total
}

func InvalidateUserListCache() {
	redis.Client.Del(redis.Ctx, UserTotalCacheKey)

	keys, err := redis.Client.SMembers(redis.Ctx, UserListKeysSet).Result()
	if err == nil && len(keys) > 0 {
		redis.Client.Del(redis.Ctx, keys...)
	}

	redis.Client.Del(redis.Ctx, UserListKeysSet)
}

func InvalidateAllSingleUserCache() {
    keys, err := redis.Client.Keys(redis.Ctx, UserCachePrefix+"*").Result()
    if err == nil && len(keys) > 0 {
        redis.Client.Del(redis.Ctx, keys...)
    }
}