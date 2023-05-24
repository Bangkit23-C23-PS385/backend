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
var dbSlave *gorm.DB

type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func dbInit(hostType string, dbName string) *gorm.DB {

	if hostType == "master" {
		hostType = os.Getenv("DB_POSTGRES_HOST_MASTER")
	} else if hostType == "slave" {
		hostType = os.Getenv("DB_POSTGRES_HOST_SLAVE")
	}

	postgresCon := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s application_name=%s",
		hostType,
		os.Getenv("DB_POSTGRES_PORT"),
		os.Getenv("DB_POSTGRES_USERNAME"),
		dbName,
		fmt.Sprintf("#%v", os.Getenv("DB_POSTGRES_PASSWORD")),
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
		dbMaster = dbInit("master", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbMaster
}

func DBSlave() *gorm.DB {
	if dbSlave == nil {
		fmt.Println("No Active Slave Connection Found")
		fmt.Println("Creating New Slave Connection")
		dbSlave = dbInit("slave", os.Getenv("DB_POSTGRES_DATABASE"))
	}
	return dbSlave
}
