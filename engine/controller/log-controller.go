package controller

import (
	"go-quantus-service/src/rabbitMQ"
	"go-quantus-service/src/redis"
	"gorm.io/gorm"
)

type LogController struct {
	Redis  *redis.RedisClient
	DB     *gorm.DB
	Rabbit *rabbitMQ.RabbitMQImpl
}

func (c *LogController) GetDependencies() *LogController {
	return c
}
