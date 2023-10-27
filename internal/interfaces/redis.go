package interfaces

import "github.com/go-redis/redis"

type Redis interface {
	IncrBy(key string, value int64) *redis.IntCmd
}
