package config

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "admin"
	dbname   = "vtt"
	port     = "5432"
	timeZone = "america/sao_paulo"
)

func initializePostgreSQL() (*gorm.DB, error) {
	logger := GetLogger("postgreSQL")

	//todo: Fazer a criação do banco de dados sozinho

	//creating database connection

	//argument to DSN
	argument := fmt.Sprintf("host= %v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v",
		host, user, password, dbname, port, timeZone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  argument,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		logger.ErrorF("postgree opening error: %v", err)
		return nil, err

	}
	//	Migrate the schema

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		logger.ErrorF("postgres  auto-migrating error: %v", err)
		return nil, err
	}

	//return db
	return db, err
}
