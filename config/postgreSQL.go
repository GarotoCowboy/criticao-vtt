package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

const (
	user     = "admin"
	password = "admin"
	dbname   = "vtt"
	port     = "9920"
	timeZone = "america/sao_paulo"
)

func initializePostgreSQL() (*gorm.DB, error) {
	logger := GetLogger("postgreSQL")
	//check if database exists
	dbPath := "./db/main.db"

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		logger.InfoF("database file does not exists, creating file...")
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return nil, err
		}

		file, err := os.Create(dbPath)

		if err != nil {
			return nil, err
		}

		file.Close()
	}

	//creating database connection

	//argument to DSN
	argument := fmt.Sprintf("user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v",
		user, password, dbname, port, timeZone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  argument,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		logger.ErrorF("postgree opening error: %v", err)
		return nil, err

	}

	//Migrate the schema

	//err = db.AutoMigrate(&schemas.Opening{})
	//if err != nil {
	//	logger.ErrorF("sqlite auto-migrating error: %v", err)
	//	return nil, err
	//}

	//return db
	return db, err
}
