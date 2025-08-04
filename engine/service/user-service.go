package service

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
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
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		if err := s.Repo.CreateUser(tx, req); err != nil {
			return nil, err
		}
		return req, nil
	})
	if err != nil {
		config.Logger.Println("err service", err.Error())
		return nil, err
	}
	return result.(*entities.User), nil
}
