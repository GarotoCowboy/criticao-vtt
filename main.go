package main

import (
	"github.com/GarotoCowboy/vttProject/api/router"
	"github.com/GarotoCowboy/vttProject/config"
)

var (
	logger *config.Logger
)

func main() {

	logger = config.GetLogger("main")

	//Initialize config
	err := config.InitPostgres()
	if err != nil {
		logger.ErrorF("config initialization error: %v", err)
		return
	}

	if err := config.CreateImgFolder(); err != nil {
		logger.ErrorF("Error creating image folder: %v", err)
		return
	}

	//Initialize the server
	router.Initializer()

}
