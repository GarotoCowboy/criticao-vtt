package router

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/handler/gameHandler"
	"github.com/GarotoCowboy/vttProject/api/handler/tablehandler"
	"github.com/GarotoCowboy/vttProject/api/handler/tableuserhandler"
	"github.com/GarotoCowboy/vttProject/api/handler/uploadHandler"
	"github.com/GarotoCowboy/vttProject/api/handler/userhandler"
	_ "github.com/GarotoCowboy/vttProject/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) {

	handler.InitializeHandler()

	var v1 = router.Group("api/v1")
	{
		//User Request
		v1.GET("/user", userhandler.GetUserHandler)
		v1.GET("/users", userhandler.ListUsersHandler)
		v1.POST("/user", userhandler.CreateUserHandler)
		v1.DELETE("/user", userhandler.DeleteUserHandler)
		v1.PUT("/user", userhandler.UpdateUserHandler)
		v1.POST("/user/upload", userhandler.UploadUserImg)

		//Table Requests

		v1.GET("/table", tablehandler.GetTableHandler)
		v1.GET("/tables", tablehandler.ListTablesHandler)
		v1.POST("/table", tablehandler.CreateTableHandler)
		v1.DELETE("/table", tablehandler.DeleteTableHandler)
		v1.PUT("/table", tablehandler.UpdateTableHandler)
		////v1.POST("/userDTO/upload", userhandler.UploadUserImg)
		//
		////TableUser Requests
		v1.GET("/tableUser", tableuserhandler.GetTableUserHandler)
		v1.GET("/tablesUsers", tableuserhandler.ListTableUsersHandler)
		v1.POST("/tableUser", tableuserhandler.CreateTableUserHandler)
		v1.POST("/tableUser/inviteLink", tableuserhandler.CreateTableUserByInviteLinkHandler)
		v1.DELETE("/tableUser", tableuserhandler.DeleteTableUserHandler)
		//v1.PUT("/tableUser", table.UpdateUserHandler)

		//gameService
		v1.POST("/:tableUser/rollDice", gameHandler.RollDiceHandler)

		//v1.POST("/table/character",characterhandler.CreateCharacterHandler)
	}

	//utils
	v1.POST("/util/upload/pdf", uploadHandler.UploadFilePDF)
	v1.POST("/util/upload/audio", uploadHandler.UploadFileMP3)

	//Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
