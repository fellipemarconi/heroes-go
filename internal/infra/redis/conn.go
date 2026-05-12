package redis

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ConnectRedis() {
	_ = godotenv.Load()

	url := os.Getenv("REDIS_URL")
	if url == "" {
		panic("REDIS_URL not set")
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	Client = redis.NewClient(opt)
}
