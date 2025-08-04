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

// ===================== HELPER =====================
func getMetaDataOrAbort(ctx *gin.Context, resp pkg.GinResponse) *pkg.TokenData {
	metaData := pkg.ExtractToken(ctx)
	if metaData == nil {
		_ = resp.SetHttpCode(http.StatusUnauthorized).ReplyFailed("97", "907", "meta data is nil", nil)
		return nil
	}
	return metaData
}

func bindContentOrAbort(ctx *gin.Context, resp pkg.GinResponse) (*pkg.RawContent, bool) {
	var req pkg.RawContent
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = resp.SetHttpCode(http.StatusUnprocessableEntity).ReplyFailed("99", "909", err.Error(), nil)
		return nil, false
	}
	return &req, true
}

func parseContentID(ctx *gin.Context) (int64, error) {
	idStr := ctx.Param("content_id")
	return strconv.ParseInt(idStr, 10, 64)
}

func (c *ContentControllerImpl) RegisterContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := getMetaDataOrAbort(ctx, resp)
	if metaData == nil {
		return
	}

	req, ok := bindContentOrAbort(ctx, resp)
	if !ok {
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
}

func (c *ContentControllerImpl) GetAllContentsController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := getMetaDataOrAbort(ctx, resp)
	if metaData == nil {
		return
	}

	result, err := c.services.GetAllContents(ctx, metaData.UserID)
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
}

func (c *ContentControllerImpl) GetContentByIDController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	if getMetaDataOrAbort(ctx, resp) == nil {
		return
	}

	id, err := parseContentID(ctx)
	if err != nil {
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "908", "invalid content_id", nil)
		return
	}

	result, err := c.services.GetContentByID(ctx, id)
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "ok", result)
}

func (c *ContentControllerImpl) UpdateContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := getMetaDataOrAbort(ctx, resp)
	if metaData == nil {
		return
	}

	req, ok := bindContentOrAbort(ctx, resp)
	if !ok {
		return
	}

	id, err := parseContentID(ctx)
	if err != nil {
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "908", "invalid content_id", nil)
		return
	}

	result, err := c.services.UpdateContent(ctx, &entities.Content{
		ID:     id,
		UserID: metaData.UserID,
		Title:  req.Title,
		Body:   req.Body,
	})
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusAccepted).ReplySuccess("00", "00001", "updated", result)
}

func (c *ContentControllerImpl) DeleteContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	if getMetaDataOrAbort(ctx, resp) == nil {
		return
	}

	id, err := parseContentID(ctx)
	if err != nil {
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "908", "invalid content_id", nil)
		return
	}

	result, err := c.services.DeleteContent(ctx, id)
	if err != nil {
		config.Logger.Println("[err service]", err.Error())
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", "deleted", result)
}

func (c *ContentControllerImpl) UpdateOrDeleteContentController(ctx *gin.Context) {
	resp := pkg.PlugGinResponse(ctx)
	metaData := getMetaDataOrAbort(ctx, resp)
	if metaData == nil {
		return
	}

	id, err := parseContentID(ctx)
	if err != nil {
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "908", "invalid content_id", nil)
		return
	}

	content := &entities.Content{
		ID:     id,
		UserID: metaData.UserID,
	}

	if ctx.Request.Method == http.MethodPut || ctx.Request.Method == http.MethodPatch {
		req, ok := bindContentOrAbort(ctx, resp)
		if !ok {
			return
		}
		content.Title = req.Title
		content.Body = req.Body
	}

	result, err := c.services.HandleContentUpdateOrDelete(ctx, content)
	if err != nil {
		_ = resp.SetHttpCode(http.StatusBadRequest).ReplyFailed("97", "90", err.Error(), nil)
		return
	}

	msg := "updated"
	if ctx.Request.Method == http.MethodDelete {
		msg = "deleted"
	}
	_ = resp.SetHttpCode(http.StatusOK).ReplySuccess("00", "00001", msg, result)
}
