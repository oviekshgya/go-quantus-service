package service

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/rabbitMQ"
	"go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(c *gin.Context, req *entities.User) (*entities.User, error)
	LoginUserController(c *gin.Context, req *entities.User) (*entities.User, error)
}

func NewUserService(db *gorm.DB, repo repository.RepositoryUser, rbt *rabbitMQ.RabbitMQImpl, rds *redis.RedisClient) UserService {
	return &UserServiceImpl{
		DB:     db,
		Repo:   repo,
		Rabbit: rbt,
		Redis:  rds,
	}
}
