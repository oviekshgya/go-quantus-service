package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
)

type UserController interface {
	RegisterUserController(c *gin.Context)
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}
