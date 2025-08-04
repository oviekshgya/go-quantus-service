package route

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-quantus-service/engine/controller"
	"go-quantus-service/engine/middleware"
	"net/http"
)

func LogInRouter(rg *gin.RouterGroup, c controller.UserController) {
	basicAuth := middleware.BasicAuthMiddleware(viper.GetString("SERVICE_USERNAME"), viper.GetString("SERVICE_PASSWORD"))
	jwtAuth := middleware.JWTAuthMiddleware()
	rg.POST("/", basicAuth, c.RegisterUserController)
	rg.POST("/login", basicAuth, c.LoginUserController)
	rg.GET("/", jwtAuth, c.UserDetailController)
	rg.GET("/:user_id", basicAuth, c.UserDetailByIDController)
	rg.PUT("/:user_id", jwtAuth, c.UpdateUserController)
	rg.DELETE("/:user_id", jwtAuth, c.DeleteUserController)
}

type InitialController struct {
	UserController controller.UserController
}

func (ctrl *InitialController) RegisterGinRoutes(router *gin.Engine) {

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	LogInRouter(router.Group("/users"), ctrl.UserController)
}
