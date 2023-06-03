package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global Variable
var dbMaster *gorm.DB

type DB struct {
	Master *gorm.DB
}

func dbInit(dbName string) *gorm.DB {
	postgresCon := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s application_name=%s",
		os.Getenv("DB_POSTGRES_HOST_MASTER"),
		os.Getenv("DB_POSTGRES_PORT"),
		os.Getenv("DB_POSTGRES_USERNAME"),
		dbName,
		os.Getenv("DB_POSTGRES_PASSWORD"),
		os.Getenv("APPLICATION_NAME"),
	)

	DB, err := gorm.Open(postgres.Open(postgresCon), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		QueryFields: true,
	})

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed connected to database")
	}

	fmt.Println("Successfully connected to database")

	return DB
}

func DBMaster() *gorm.DB {
	if dbMaster == nil {
		fmt.Println("No Active Master Connection Found")
		fmt.Println("Creating New Master Connection")
		dbMaster = dbInit(os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbMaster
}
