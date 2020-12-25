package storage

import (
	"brank/core"
	"net/url"

	"github.com/go-redis/redis"
)

func NewRedis(config *core.Config) *redis.Client {
	env := core.GetEnvironment()
	if env == core.Staging || env == core.Production {
		parsedURL, _ := url.Parse(config.REDIS_URL)
		password, _ := parsedURL.User.Password()
		return redis.NewClient(&redis.Options{
			Addr:     parsedURL.Host,
			Password: password,
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDRESS,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})
}
