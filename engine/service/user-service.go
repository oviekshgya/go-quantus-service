package service

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/rabbitMQ"
	rds "go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB     *gorm.DB
	Repo   repository.RepositoryUser
	Rabbit *rabbitMQ.RabbitMQImpl
	Redis  *rds.RedisClient
}

func (s *UserServiceImpl) RegisterUser(c *gin.Context, req *entities.User) (*entities.User, error) {
	return &entities.User{}, nil
}
