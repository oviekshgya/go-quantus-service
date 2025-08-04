package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
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
}

func NewContentController(userService service.ContentSerivce) ContentController {
	return &ContentControllerImpl{
		services: userService,
	}
}
