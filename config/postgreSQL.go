package config

import (
	"fmt"
	"os"

	"github.com/GarotoCowboy/vttProject/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbname   = "vtt"
	port     = "5432"
	timeZone = "america/Sao_Paulo"
)

func createDataBaseIfNotExists(logger *Logger) error {
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")

	//Connection with default database "postgres"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database %v", err)
	}

	//Close connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to obtain sql.DB: %v", err)
	}
	defer sqlDB.Close()

	//Verify if database exists
	var exists bool
	err = db.Raw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", dbname).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database %v exists %v", dbname, err)
	}

	//Creating Database
	if !exists {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", dbname)).Error
		if err != nil {
			return fmt.Errorf("failed to create database %v", dbname)
		}
		logger.InfoF("Database %v created", dbname)
	} else {
		logger.InfoF("Database %v already exists", dbname)
	}

	return nil
}

func initializePostgreSQL() (*gorm.DB, error) {

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")

	logger := GetLogger("postgreSQL")

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
	//err = db.Migrator().DropTable(&models.TableUser{})
	err = db.AutoMigrate(
		&models.User{},
		&models.Table{},
		&models.TableUser{},
		&models.Character{},
		&models.ChatMessage{},
		&models.Scene{},
		&models.Image{},
		&models.GameObjectOwner{},
		&models.PlacedImage{},
		&models.PlacedToken{},
		&models.Token{},
		&models.Bar{})
	if err != nil {
		logger.ErrorF("postgres  auto-migrating error: %v", err)
		return nil, err
	}

	//return db
	return db, err
}
