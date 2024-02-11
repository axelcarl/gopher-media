package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func InitRedis(address, password, dbStr string) (*redis.Client, error) {
	database, err := strconv.Atoi(dbStr)

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       database,
	})

	return rdb, nil
}

func CreateSession(rdb *redis.Client, userID string, expiry time.Duration) (string, error) {
	sessionID := uuid.New().String()

	err := rdb.Set(ctx, sessionID, userID, expiry).Err()

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func GetUserIDBySession(rdb *redis.Client, sessionID string) (string, error) {
	userID, err := rdb.Get(ctx, sessionID).Result()

	if err != nil {
		return "", err
	}

	return userID, nil
}

func DeleteSession(rdb *redis.Client, sessionId string) error {
	return rdb.Del(ctx, sessionId).Err()
}
