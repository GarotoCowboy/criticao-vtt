package config

import (
	"fmt"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
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
