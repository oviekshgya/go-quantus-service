package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
	"net/http"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func (c *UserControllerImpl) RegisterUserController(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": "success",
	})
	return
}
