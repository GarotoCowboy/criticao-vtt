package main

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/server"
	"github.com/GarotoCowboy/vttProject/api/router"
	"github.com/GarotoCowboy/vttProject/config"
)

var (
	logger *config.Logger
)

// @title VTT API
// @version 1.0.0
// @description The VTT API for managing a VTT.
// @contact.name Pedro Henrique Marques
// @contact.email comercial.pedromarques@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {

	logger = config.GetLogger("main")

	//Initialize config
	err := config.InitPostgres()
	if err != nil {
		logger.ErrorF("config initialization error: %v", err)
		return
	}
	//create img folder if not exits
	if err := config.CreateImgFolder(); err != nil {
		logger.ErrorF("Error... Creating image folder: %v", err)
		return
	}

	//create a file folder if not exists example pdf
	if err := config.CreateFileFolder(); err != nil {
		logger.ErrorF("Error... Creating file folder: %v", err)
		return
	}

	db := config.GetPostgreSQL()

	//Initialize the server
	go server.RunGRPCServer(db, logger)
	router.Initializer()

}
