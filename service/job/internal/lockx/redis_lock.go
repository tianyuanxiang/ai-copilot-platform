package lockx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var releaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)

func TryAcquire(ctx context.Context, rdb *redis.Client, key, owner string, ttl time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, owner, ttl).Result()
}

func Release(ctx context.Context, rdb *redis.Client, key, owner string) error {
	return releaseScript.Run(ctx, rdb, []string{key}, owner).Err()
}
