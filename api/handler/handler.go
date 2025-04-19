package handler

import (
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

var (
	Logger *config.Logger
	db     *gorm.DB
)

func InitializeHandler() {
	Logger = config.GetLogger("handler")
	db = config.GetPostgreSQL()
}
