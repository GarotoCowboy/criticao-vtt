package handler

import (
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db = config.GetPostgreSQL()
}

func GetHandlerLogger() *config.Logger {
	return logger
}

func GetHandlerDB() *gorm.DB {
	return db
}
