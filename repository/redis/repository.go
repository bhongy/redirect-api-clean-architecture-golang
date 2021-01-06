package redis

import (
	"fmt"
	"time"

	"github.com/bhongy/tmp-clean-arch-golang/shortener"
	"github.com/go-redis/redis"
)

var (
	TimeFormat = time.RFC3339
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, fmt.Errorf("repository.NewRedisRepository: %v", err)
	}
	return &redisRepository{client}, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("repository.Redirect.Find: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("repository.Redirect.Find: %w", shortener.ErrRedirectNotFound)
	}

	createdAt, err := time.Parse(TimeFormat, data["created_at"])
	if err != nil {
		return nil, fmt.Errorf("repository.Redirect.Find: %v", err)
	}

	redirect := shortener.Redirect{
		Code:      data["code"],
		URL:       data["url"],
		CreatedAt: createdAt,
	}
	return &redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt.Format(TimeFormat),
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return fmt.Errorf("repository.Redirect.Store: %v", err)
	}
	return nil
}
