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

func GetLogger() *config.Logger {
	return logger
}

func GetDB() *gorm.DB {
	return db
}
