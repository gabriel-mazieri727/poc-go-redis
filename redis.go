package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect to Redis
func (r *Redis) Connect() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
	return redisClient, nil
}

// Close the Redis connection
func (r *Redis) Close(redisClient *redis.Client) error {
	return redisClient.Close()
}

// Set a key-value pair in Redis
func (r *Redis) Set(redisClient *redis.Client, key string, value string) error {
	return redisClient.Set(context.Background(), key, value, 0).Err()
}

// Set a key-value pair in Redis with a TTL (Time To Live || Expiration Time)
func (r *Redis) SetWithTTL(redisClient *redis.Client, key string, value string, ttl time.Duration) error {
	return redisClient.Set(context.Background(), key, value, ttl).Err()
}

// Get a value from Redis
func (r *Redis) Get(redisClient *redis.Client, key string) (string, error) {
	value, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// Delete a key-value pair from Redis
func (r *Redis) Delete(redisClient *redis.Client, key string) error {
	return redisClient.Del(context.Background(), key).Err()
}
