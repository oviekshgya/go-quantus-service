package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"net/http"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func (c *UserControllerImpl) RegisterUserController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	var req pkg.RawUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("99", "909", err.Error(), nil)

		return
	}
	result, err := c.UserService.RegisterUser(ctx, &entities.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		config.Logger.Println("err service", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("99", "909", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusCreated).ReplySuccess("00", "00001", "created", result)
	return
}
