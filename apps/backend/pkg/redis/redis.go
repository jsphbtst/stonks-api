package cache

import (
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

var cache = &Cache{}

func Init(uri string) (*redis.Client, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	cache.Client = client

	return client, nil
}
