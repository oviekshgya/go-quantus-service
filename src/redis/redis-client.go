package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go-quantus-service/src/entities"
	"time"
)

type RedisClient struct {
	C *redis.Client
}

type RedisConfig interface {
	GetKey(key string, src interface{}) error
	SetKey(key string, value interface{}, expiration time.Duration) error
	DeleteKey(key string) error
	SettexKey(key string, value interface{}, expiration time.Duration) error
	ExpireKey(key string, expiration time.Duration) error
	FlushAll() error
	Close() error
	PushLogToQueue(queueName string, logData interface{}) error
	PopLogsFromQueue(queueName string, batchSize int) ([]entities.LogEntry, error)
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

func (client *RedisClient) PushLogToQueue(queueName string, logData interface{}) error {
	data, err := json.Marshal(logData)
	if err != nil {
		return err
	}
	return client.C.LPush(queueName, data).Err()
}

// PopLogsFromQueue pops multiple logs from Redis list (queue), up to batchSize
func (client *RedisClient) PopLogsFromQueue(queueName string, batchSize int) ([]entities.LogEntry, error) {
	var logs []entities.LogEntry

	for i := 0; i < batchSize; i++ {
		data, err := client.C.RPop(queueName).Result()
		if err == redis.Nil {
			break // queue empty
		}
		if err != nil {
			return logs, err
		}
		var log entities.LogEntry
		if err := json.Unmarshal([]byte(data), &log); err == nil {
			logs = append(logs, log)
		}
	}
	return logs, nil
}
