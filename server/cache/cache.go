// cache/cache.go
package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// InitCache initializes the Redis client with given parameters and handles connection errors
func InitCache(addr string, password string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // leave empty if no password
		DB:       db,       // use default DB
	})
	// Immediate ping to check if Redis server is accessible
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Could not connect to Redis: %v", err)
		return errors.New("failed to connect to Redis")
	}
	log.Println("Connected to Redis successfully")
	return nil
	// // Test connection
	// if err := rdb.Ping(ctx).Err(); err != nil {
	//     log.Printf("Could not connect to Redis: %v", err)
	//     return errors.New("failed to connect to Redis")
	// }
	// log.Println("Connected to Redis successfully")
	// return nil
}

// SetCacheValue sets a key-value pair in Redis with error handling
func SetCacheValue(key string, value interface{}) error {
	if err := rdb.Set(ctx, key, value, 0).Err(); err != nil {
		log.Printf("Failed to set cache value: %v", err)
		return errors.New("failed to set cache value")
	}
	return nil
}

// GetCacheValue retrieves a value from Redis by key, with error handling
func GetCacheValue(key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Cache miss for key: %s", key)
		return "", errors.New("cache miss")
	} else if err != nil {
		log.Printf("Failed to get cache value: %v", err)
		return "", errors.New("failed to get cache value")
	}
	return value, nil
}

// DeleteCacheKey deletes a key in Redis with error handling
func DeleteCacheKey(key string) error {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("Failed to delete cache key: %v", err)
		return errors.New("failed to delete cache key")
	}
	return nil
}
