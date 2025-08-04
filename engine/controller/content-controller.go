package controller

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/engine/service"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"net/http"
	"strconv"
)

type ContentControllerImpl struct {
	services service.ContentSerivce
}

func (c *ContentControllerImpl) RegisterContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}

	var req pkg.RawContent
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("99", "909", err.Error(), nil)
		return
	}

	result, err := c.services.RegisterContent(ctx, &entities.Content{
		UserID: metaData.UserID,
		Title:  req.Title,
		Body:   req.Body,
	})
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusCreated).ReplySuccess("00", "00001", "created", result)
	return
}

func (c *ContentControllerImpl) GetAllContentsController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}

	result, err := c.services.GetAllContents(ctx, metaData.UserID)
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
	return
}

func (c *ContentControllerImpl) GetContentByIDController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}
	id, _ := strconv.Atoi(ctx.Param("content_id"))

	result, err := c.services.GetContentByID(ctx, int64(id))
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
	return
}

func (c *ContentControllerImpl) UpdateContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}

	var req pkg.RawContent
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("99", "909", err.Error(), nil)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("content_id"))
	result, err := c.services.UpdateContent(ctx, &entities.Content{
		UserID: metaData.UserID,
		Title:  req.Title,
		Body:   req.Body,
		ID:     int64(id),
	})
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusAccepted).ReplySuccess("00", "00001", "updated", result)
	return
}

func (c *ContentControllerImpl) DeleteContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return
	}
	id, _ := strconv.Atoi(ctx.Param("content_id"))
	result, err := c.services.DeleteContent(ctx, int64(id))
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "updated", result)
	return
}
