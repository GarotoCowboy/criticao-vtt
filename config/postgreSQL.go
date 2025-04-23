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
	timeZone = "america/Sao_Paulo"
)

func createDataBaseIfNotExists(logger *Logger) error {

	//Connection with default database "postgres"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connect to database %v", err)
	}

	//Close connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed to obtain sql.DB: %v", err)
	}
	defer sqlDB.Close()

	//Verify if database exists
	var exists bool
	err = db.Raw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", dbname).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("Failed to check if database %v exists %v", dbname, err)
	}

	//Creating Database
	if !exists {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", dbname)).Error
		if err != nil {
			return fmt.Errorf("Failed to create database %v", dbname)
		}
		logger.InfoF("Database %v created", dbname)
	} else {
		logger.InfoF("Database %v already exists", dbname)
	}

	return nil
}

func initializePostgreSQL() (*gorm.DB, error) {
	logger := GetLogger("postgreSQL")

	//todo: Fazer a criação do banco de dados sozinho
	if err := createDataBaseIfNotExists(logger); err != nil {
		logger.ErrorF("Failed to create database %v", err)
		return nil, fmt.Errorf("error to initializate PostgreSQL: %v", err)
	}

	//creating database connection

	//argument to DSN
	argument := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v",
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
