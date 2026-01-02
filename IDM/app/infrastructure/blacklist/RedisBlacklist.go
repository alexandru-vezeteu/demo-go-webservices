package blacklist

import (
	"context"
	"time"

	"idmService/application/domain"

	"github.com/redis/go-redis/v9"
)

const (
	blacklistKeyPrefix = "token:blacklist:"
	defaultTTL         = 24 * time.Hour
)

type RedisBlacklist struct {
	client *redis.Client
	ttl    time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
}

func NewRedisBlacklist(config *RedisConfig) (*RedisBlacklist, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, &domain.InternalError{Operation: "redis connection", Err: err}
	}

	ttl := config.TTL
	if ttl == 0 {
		ttl = defaultTTL
	}

	return &RedisBlacklist{
		client: client,
		ttl:    ttl,
	}, nil
}

func (b *RedisBlacklist) Add(token string, reason string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = time.Minute
	}

	key := blacklistKeyPrefix + token
	err := b.client.Set(ctx, key, reason, ttl).Err()
	if err != nil {
		return &domain.InternalError{Operation: "redis set", Err: err}
	}

	return nil
}

func (b *RedisBlacklist) IsBlacklisted(token string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := blacklistKeyPrefix + token
	reason, err := b.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, ""
	}
	if err != nil {
		return false, ""
	}

	return true, reason
}

func (b *RedisBlacklist) Remove(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := blacklistKeyPrefix + token
	result, err := b.client.Del(ctx, key).Result()
	if err != nil {
		return &domain.InternalError{Operation: "redis delete", Err: err}
	}
	if result == 0 {
		return &domain.NotFoundError{Resource: "token", ID: ""}
	}

	return nil
}

func (b *RedisBlacklist) Close() error {
	return b.client.Close()
}
