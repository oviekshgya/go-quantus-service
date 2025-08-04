package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"time"
)

type RedisClient struct {
	C *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("REDIS_ADDR"),
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     50,
		MinIdleConns: 10,
		IdleTimeout:  300 * time.Second,
	})

	if err := client.Ping().Err(); err != nil {
		fmt.Println("Unable to connect to redis:", err)
	}

	fmt.Println("REDIS CONNECTED")
	return &RedisClient{C: client}
}

type LoggerRedis struct {
	Code         string      `json:"code"`
	Timestamp    time.Time   `json:"timestamp"`
	Id           int         `json:"id"`
	Repositories string      `json:"repositories"`
	Column       int         `json:"column"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func (client *RedisClient) GetKey(key string, src interface{}) error {
	val, err := client.C.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), &src)
}

func (client *RedisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.C.Set(key, cacheEntry, expiration).Err()
}

func (client *RedisClient) DeleteKey(key string) error {
	iter := client.C.Scan(0, key, 0).Iterator()
	for iter.Next() {
		if err := client.C.Del(iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (client *RedisClient) SettexKey(key string, value interface{}, expiration time.Duration) error {
	cacheSettex, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.C.SetXX(key, cacheSettex, expiration).Err()
}

func (client *RedisClient) ExpireKey(key string, expiration time.Duration) error {
	return client.C.Expire(key, expiration).Err()
}

func (client *RedisClient) FlushAll() error {
	return client.C.FlushAll().Err()
}

func (client *RedisClient) Close() error {
	return client.C.Close()
}
