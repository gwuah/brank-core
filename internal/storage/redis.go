package storage

import (
	"brank/internal"
	"net/url"

	"github.com/go-redis/redis"
)

func NewRedis(config *internal.Config) *redis.Client {
	env := internal.GetEnvironment()
	if env == internal.Staging || env == internal.Production {
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
