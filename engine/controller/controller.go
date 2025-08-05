package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
	"go-quantus-service/src/rabbitMQ"
	"go-quantus-service/src/redis"
	"gorm.io/gorm"
)

type UserController interface {
	RegisterUserController(c *gin.Context)
	LoginUserController(c *gin.Context)
	UserDetailController(c *gin.Context)
	UserDetailByIDController(c *gin.Context)
	UpdateUserController(c *gin.Context)
	DeleteUserController(c *gin.Context)
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

type ContentController interface {
	RegisterContentController(c *gin.Context)
	GetAllContentsController(c *gin.Context)
	GetContentByIDController(c *gin.Context)
	UpdateContentController(c *gin.Context)
	DeleteContentController(c *gin.Context)
	UpdateOrDeleteContentController(ctx *gin.Context)
}

func NewContentController(userService service.ContentSerivce) ContentController {
	return &ContentControllerImpl{
		services: userService,
	}
}

type LogControllerinterface interface {
	GetDependencies() *LogController
}

func NewLogController(db *gorm.DB, redisClient *redis.RedisClient, rabbit *rabbitMQ.RabbitMQImpl) LogControllerinterface {
	return &LogController{
		Redis:  redisClient,
		DB:     db,
		Rabbit: rabbit,
	}
}
