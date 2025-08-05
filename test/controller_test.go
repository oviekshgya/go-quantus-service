package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-quantus-service/engine/controller"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"go-quantus-service/test/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/", nil)
	c.Request = req
	return c
}

//func TestRegisterUserController(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockCtrl := mocks.NewMockUserController(ctrl)
//	ctx := setupTestContext()
//
//	mockCtrl.EXPECT().RegisterUserController(ctx).Times(1)
//
//	mockCtrl.RegisterUserController(ctx)
//}

func TestRegisterUserController_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock service
	mockService := mocks.NewMockUserService(ctrl)

	// Inisialisasi controller dengan mock service
	userCtrl := &controller.UserControllerImpl{
		UserService: mockService,
	}

	// Test request payload
	reqBody := pkg.RawUser{
		FullName: "John Doe",
		Email:    "john@example.com",
		Password: "123456",
		Role:     "user",
		IsActive: true,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Simulasi HTTP request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Mock expected service call
	expectedUser := &entities.User{
		FullName: "John Doe",
		Email:    "john@example.com",
		Password: "123456",
		Role:     "user",
		IsActive: true,
	}
	mockService.EXPECT().RegisterUser(gomock.Any(), expectedUser).Return(expectedUser, nil)

	// Jalankan controller
	userCtrl.RegisterUserController(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"status_number":"00"`)
	assert.Contains(t, w.Body.String(), `"status_message":"created"`)
}

func TestLoginUserController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	c := controller.UserControllerImpl{UserService: mockService}

	router := gin.Default()
	router.POST("/login", c.LoginUserController)

	body := pkg.RawLogin{Email: "test@example.com", Password: "secret"}
	mockService.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(&entities.User{ID: 1, Role: "user"}, nil)

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserDetailByIDController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	c := controller.UserControllerImpl{UserService: mockService}

	router := gin.Default()
	router.GET("/user/:user_id", c.UserDetailByIDController)

	mockService.EXPECT().UserDetail(gomock.Any(), int64(123)).Return(&entities.User{ID: 123, Email: "test@example.com"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/user/123", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteUserController_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	c := controller.UserControllerImpl{UserService: mockService}

	router := gin.Default()
	router.DELETE("/user/:user_id", c.DeleteUserController)

	mockService.EXPECT().DeleteUser(gomock.Any(), int64(10)).Return(nil, errors.New("delete error"))

	req, _ := http.NewRequest(http.MethodDelete, "/user/10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
