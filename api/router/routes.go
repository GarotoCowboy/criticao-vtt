package router

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/handler/userhandler"
	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {

	handler.InitializeHandler()

	var v1 = router.Group("api/v1")
	{

		v1.GET("/user", userhandler.GetUserHandler)
		v1.GET("/users", userhandler.ListUsersHandler)
		v1.POST("/user", userhandler.CreateUserHandler)
		v1.DELETE("/user", userhandler.DeleteUserHandler)
		v1.PUT("/user")

	}

}
