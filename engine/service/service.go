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
	LoginUser(c *gin.Context, req *entities.User) (*entities.User, error)
	UserDetail(c *gin.Context, id int64) (*entities.User, error)
	UpdateUSser(c *gin.Context, req *entities.User) (*entities.User, error)
	DeleteUser(c *gin.Context, id int64) (*int64, error)
}

func NewUserService(db *gorm.DB, repo repository.RepositoryUser, rbt *rabbitMQ.RabbitMQImpl, rds *redis.RedisClient) UserService {
	return &UserServiceImpl{
		DB:     db,
		Repo:   repo,
		Rabbit: rbt,
		Redis:  rds,
	}
}

type ContentSerivce interface {
	RegisterContent(c *gin.Context, req *entities.Content) (*entities.Content, error)
	GetAllContents(c *gin.Context, id int64) ([]entities.Content, error)
	GetContentByID(c *gin.Context, id int64) (*entities.Content, error)
}

func NewContentService(db *gorm.DB, repo repository.RepositoryContent, rbt *rabbitMQ.RabbitMQImpl, rds *redis.RedisClient) ContentSerivce {
	return &ContentServiceImpl{
		DB:     db,
		Repo:   repo,
		Rabbit: rbt,
		Redis:  rds,
	}
}
