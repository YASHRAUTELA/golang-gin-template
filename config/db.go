package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func DBConfig() string {
	db_Host := os.Getenv("DB_HOST")
	db_User := os.Getenv("DB_USER")
	db_Password := os.Getenv("DB_PASSWORD")
	db_Name := os.Getenv("DB_NAME")
	db_Port := os.Getenv("DB_PORT")
	db_Timezone := os.Getenv("DB_TZ")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", db_Host, db_User, db_Password, db_Name, db_Port, db_Timezone)
	return dsn
}

func InitDBConnection() {
	dsn := DBConfig()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("database connection failed")
		// Add logger
	}
}
