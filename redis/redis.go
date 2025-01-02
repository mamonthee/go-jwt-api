package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Update with your Redis configuration
})

// AddTokenToBlacklist adds a token to the blacklist with an expiration time.
func AddTokenToBlacklist(token string, expiry time.Duration) error {
	return RedisClient.Set(ctx, token, "blacklisted", expiry).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted.
func IsTokenBlacklisted(token string) bool {
	_, err := RedisClient.Get(ctx, token).Result()
	return err == nil // If no error, the token is blacklisted
}
