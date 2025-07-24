package config

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

var (
	db     *gorm.DB
	logger *Logger
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

	if _, err := os.Stat("./img"); os.IsNotExist(err) {
		logger.InfoF("img folder does not exist, creating...")
		if err := os.Mkdir("./img", os.ModeDir); err != nil {
			return fmt.Errorf("Error creating image folder: %v", err)
		}
		logger.InfoF("Create img folder successfully")
	}

	return nil
}

func CreateFileFolder() error {

	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		logger.InfoF("files folder does not exist, creating...")
		if err := os.Mkdir("./files", os.ModeDir); err != nil {
			return fmt.Errorf("Error creating files folder: %v", err)
		}
		logger.InfoF("Create files folder successfully")
	}

	return nil
}
