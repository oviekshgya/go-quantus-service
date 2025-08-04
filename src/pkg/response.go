package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinResponse interface {
	SetHttpCode(code int) GinResponse
	ReplyAs(res GinResponse) error
	Reply(status int, number string, code string, message string, data ...any) error
	ReplyFailed(number string, code string, message string, data ...any) error
	ReplySuccess(number string, code string, message string, data ...any) error
	ReplyCustom(httpStatusCode int, res any) error
	HttpStatusCode() int
	SetHttpStatusCode(httpStatusCode int) GinResponse
	GetStatus() int
	GetStatusNumber() string
	GetStatusCode() string
	GetStatusMessage() string
	GetData() any
}

type GinResponseX struct {
	c              *gin.Context
	httpStatusCode int
	Status         int    `json:"status"`
	StatusNumber   string `json:"status_number"`
	StatusCode     string `json:"status_code"`
	StatusMessage  string `json:"status_message"`
	Data           any    `json:"data"`
}

func PlugGinResponse(c *gin.Context) GinResponse {
	return &GinResponseX{
		c: c,
	}
}

func (r *GinResponseX) HttpStatusCode() int {
	return r.httpStatusCode
}

func (r *GinResponseX) SetHttpStatusCode(httpStatusCode int) GinResponse {
	r.httpStatusCode = httpStatusCode
	return r
}

func (r *GinResponseX) GetStatus() int {
	return r.Status
}

func (r *GinResponseX) GetStatusNumber() string {
	return r.StatusNumber
}

func (r *GinResponseX) GetStatusCode() string {
	return r.StatusCode
}

func (r *GinResponseX) GetStatusMessage() string {
	return r.StatusMessage
}

func (r *GinResponseX) GetData() any {
	return r.Data
}

func (r *GinResponseX) SetHttpCode(code int) GinResponse {
	r.httpStatusCode = code
	return r
}

func (r *GinResponseX) ReplyAs(res GinResponse) error {
	r.Status = res.GetStatus()
	r.StatusNumber = res.GetStatusNumber()
	r.StatusCode = res.GetStatusCode()
	r.StatusMessage = res.GetStatusMessage()
	r.Data = res.GetData()

	r.c.JSON(res.HttpStatusCode(), r)
	return nil
}

func (r *GinResponseX) Reply(status int, number string, code string, message string, data ...any) error {
	r.Status = status
	r.StatusNumber = number
	r.StatusCode = code
	r.StatusMessage = message
	if len(data) > 0 {
		r.Data = data[0]
	}
	r.c.JSON(r.httpStatusCodeOrDefault(), r)
	return nil
}

func (r *GinResponseX) ReplyFailed(number string, code string, message string, data ...any) error {
	r.Status = 0
	r.httpStatusCode = http.StatusBadRequest
	return r.Reply(0, number, code, message, data...)
}

func (r *GinResponseX) ReplySuccess(number string, code string, message string, data ...any) error {
	r.Status = 1
	r.httpStatusCode = http.StatusOK
	return r.Reply(1, number, code, message, data...)
}

func (r *GinResponseX) ReplyCustom(httpStatusCode int, res any) error {
	r.c.JSON(httpStatusCode, res)
	return nil
}

func (r *GinResponseX) httpStatusCodeOrDefault() int {
	if r.httpStatusCode == 0 {
		return http.StatusOK
	}
	return r.httpStatusCode
}
