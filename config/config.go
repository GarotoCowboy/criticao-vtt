package config

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	logger     *Logger
	JWT_SECRET []byte
)

func InitPostgres() error {
	var err error

	//Initialize Postgrees
	db, err = initializePostgreSQL()
	if err != nil {
		return fmt.Errorf("error initializing PostgreSQL: %v", err)
	}

	return nil
}

func GetPostgreSQL() *gorm.DB { return db }

func GetLogger(p string) *Logger {
	logger = newLogger(p)
	return logger
}

func CreateImgFolder() error {

	if _, err := os.Stat("vttData/libraryImg"); os.IsNotExist(err) {
		logger.InfoF("libraryImg folder does not exist, creating...")
		if err := os.MkdirAll("vttData/libraryImg", os.ModeDir); err != nil {
			return fmt.Errorf("Error creating image folder: %v", err)
		}
		logger.InfoF("Create libraryImg folder successfully")
	}

	return nil
}

func CreateFileFolder() error {

	if _, err := os.Stat("vttData/files"); os.IsNotExist(err) {
		logger.InfoF("files folder does not exist, creating...")
		if err := os.MkdirAll("vttData/files", os.ModeDir); err != nil {
			return fmt.Errorf("Error creating files folder: %v", err)
		}
		logger.InfoF("Create files folder successfully")
	}

	return nil
}
