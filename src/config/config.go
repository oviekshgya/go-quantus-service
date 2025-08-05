package config

import (
	"encoding/json"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/redis"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var NRc *newrelic.Application

var Logger *logrus.Logger
var LoggerEntry *logrus.Entry

func init() {
	pwd, _ := os.Getwd()
	viper.SetConfigFile(fmt.Sprintf("%s/.env", pwd))
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error saat membaca file .env: %v", err)
	}

	name := viper.GetString("LOG_NAME")
	f, err := os.OpenFile(fmt.Sprintf("%s.log", name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(f)
	LoggerEntry = logrus.NewEntry(Logger)

	NRc, _ = newrelic.NewApplication(
		newrelic.ConfigAppName(viper.GetString("NEW_RELIC_APP_NAME")),
		newrelic.ConfigLicense(viper.GetString("NEW_RELIC_LICENSE")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

}

func StartLogWorker(db *gorm.DB, redisClient *redis.RedisClient, batchSize int, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			queueLen, err := redisClient.C.LLen("log_queue").Result()
			if err != nil {
				fmt.Println("Redis LLEN error:", err)
				continue
			}
			if queueLen >= int64(batchSize) {
				logStrings, err := redisClient.C.LRange("log_queue", 0, int64(batchSize-1)).Result()
				if err != nil {
					fmt.Println("Redis LRANGE error:", err)
					continue
				}

				var logs []entities.LogEntry
				for _, item := range logStrings {
					var log entities.LogEntry
					if err := json.Unmarshal([]byte(item), &log); err == nil {
						logs = append(logs, log)
					}
				}

				if len(logs) > 0 {
					if err := db.Create(&logs).Error; err != nil {
						fmt.Println("DB insert error:", err)
						continue
					}

					if err := redisClient.C.LTrim("log_queue", int64(batchSize), -1).Err(); err != nil {
						fmt.Println("Redis LTRIM error:", err)
					}

				}
			}
		}
	}()
}
