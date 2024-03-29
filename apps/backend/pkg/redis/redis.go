package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	Ctx    context.Context
}

var cache = &Cache{}

func Init(uri string, ctx context.Context) (*redis.Client, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	cache.Client = client
	cache.Ctx = ctx

	return client, nil
}
