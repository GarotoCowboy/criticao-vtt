package router

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/handler/authhandler"
	"github.com/GarotoCowboy/vttProject/api/handler/gameHandler"
	"github.com/GarotoCowboy/vttProject/api/handler/tablehandler"
	"github.com/GarotoCowboy/vttProject/api/handler/tableuserhandler"
	"github.com/GarotoCowboy/vttProject/api/handler/uploadHandler"
	"github.com/GarotoCowboy/vttProject/api/handler/userhandler"
	"github.com/GarotoCowboy/vttProject/api/middleware"
	_ "github.com/GarotoCowboy/vttProject/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) {

	handler.InitializeHandler()

	var v1 = router.Group("api/v1")
	{
		//todo: Public routes
		v1.POST("/user", userhandler.CreateUserHandler)
		v1.POST("/login", authhandler.LoginHandler)

		//todo: Authenticated routes
		authenticated := v1.Group("/")
		authenticated.Use(middleware.AuthMiddleware())
		{
			//todo:User Request
			authenticated.GET("/user", userhandler.GetUserHandler)
			authenticated.GET("/users", userhandler.ListUsersHandler)

			authenticated.DELETE("/user/me", userhandler.DeleteUserHandler)
			authenticated.PUT("/user/me", userhandler.UpdateUserHandler)
			authenticated.POST("/user/upload", userhandler.UploadUserImg)

			//todo:Table Requests

			authenticated.GET("/table", tablehandler.GetTableHandler)
			authenticated.GET("/tables", tablehandler.ListTablesHandler)
			authenticated.POST("/table", tablehandler.CreateTableHandler)
			authenticated.DELETE("/table", tablehandler.DeleteTableHandler)
			authenticated.PUT("/table", tablehandler.UpdateTableHandler)
			authenticated.POST("/table/:table_id/attachment", uploadHandler.UploadAttachmentsHandler)
			//v1.POST("/userDTO/upload", userhandler.UploadUserImg)
			//
			//todo:TableUser Requests
			authenticated.GET("/tableUser", tableuserhandler.GetTableUserHandler)
			authenticated.GET("/tablesUsers", tableuserhandler.ListTableUsersHandler)
			authenticated.POST("/tableUser", tableuserhandler.CreateTableUserHandler)
			authenticated.POST("/tableUser/inviteLink", tableuserhandler.CreateTableUserByInviteLinkHandler)
			authenticated.DELETE("/tableUser", tableuserhandler.DeleteTableUserHandler)
			//v1.PUT("/tableUser", table.UpdateUserHandler)

			//gameService
			authenticated.POST("/tables/:tableID/roll", gameHandler.RollDiceHandler)

			//v1.POST("/table/character",characterhandler.CreateCharacterHandler)
		}

	}

	//utils

	//Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
