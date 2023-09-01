package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	Client *redis.Client
}

type RedisConfig struct {
	Address  string
	Password string
}

func OpenRedis(cfg RedisConfig) Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	})
	return Cache{
		Client: redisClient,
	}
}
