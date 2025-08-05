package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"net/http"
	"strconv"
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
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("99", "909", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusCreated).Reply(http.StatusCreated, "00", "00001", "created", result)
	return
}

func (c *UserControllerImpl) LoginUserController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	var req pkg.RawLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("98", "908", err.Error(), nil)
		return
	}

	result, err := c.UserService.LoginUser(ctx, &entities.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("98", "908", err.Error(), nil)
		return
	}
	genTok, _ := pkg.GenerateJWT(fmt.Sprintf("%d", result.ID), fmt.Sprintf("%s", result.Role), 30)
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", map[string]interface{}{
		"token": genTok,
	})

}

func (c *UserControllerImpl) UserDetailController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}
	result, err := c.UserService.UserDetail(ctx, metaData.UserID)
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
}

func (c *UserControllerImpl) UserDetailByIDController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	id, _ := strconv.Atoi(ctx.Param("user_id"))
	result, err := c.UserService.UserDetail(ctx, int64(id))
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
}

func (c *UserControllerImpl) UpdateUserController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}
	var raw pkg.RawUser
	if err := ctx.ShouldBindJSON(&raw); err != nil {
		config.Logger.Println("[err validation]", err.Error())
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("99", "909", err.Error(), nil)
		return
	}
	id, _ := strconv.Atoi(ctx.Param("user_id"))
	result, er := c.UserService.UpdateUSser(ctx, &entities.User{
		ID:       int64(id),
		FullName: raw.FullName,
		Email:    raw.Email,
		Password: raw.Password,
		Role:     raw.Role,
		IsActive: raw.IsActive,
	})
	if er != nil {
		config.Logger.Println("[err service]", er.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("99", "909", er.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusAccepted).ReplySuccess("01", "00001", "accept", result)
	return
}

func (c *UserControllerImpl) DeleteUserController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	id, _ := strconv.Atoi(ctx.Param("user_id"))
	result, err := c.UserService.DeleteUser(ctx, int64(id))
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("02", "00002", "delete", result)
}
