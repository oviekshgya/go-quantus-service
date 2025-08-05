package app

import (
	"github.com/google/wire"
	"go-quantus-service/engine/controller"
	"go-quantus-service/engine/service"
	"go-quantus-service/src/db"
	"go-quantus-service/src/rabbitMQ"
	"go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
)

func InitializeUserControllers() (controller.UserController, error) {
	wire.Build(
		controller.NewUserController,
		service.NewUserService,
		redis.NewRedisClient,
		rabbitMQ.NewRabbitMQConnection,
		repository.NewUSerRepository,
		db.GetDB)
	return nil, nil
}

func InitializeContentControllers() (controller.ContentController, error) {
	wire.Build(
		controller.NewContentController,
		service.NewContentService,
		redis.NewRedisClient,
		rabbitMQ.NewRabbitMQConnection,
		repository.NewContentRepository,
		db.GetDB)
	return nil, nil
}

func InitializeLogControllers() (controller.LogControllerinterface, error) {
	wire.Build(
		controller.NewLogController,
		redis.NewRedisClient,
		rabbitMQ.NewRabbitMQConnection,
		db.GetDB)
	return nil, nil
}
