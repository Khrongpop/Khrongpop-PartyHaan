
package db

import (
	"main/model"
	"fmt"
	"gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB
const (
	dbHost   = "db"
	port     = 5432
	username = "root"
	password = "root"
	dbName   = "PartyHaan"
	TimeZone = "Asia/Bangkok"
)


// GetDB - call this method to get db
func GetDB() *gorm.DB {
	return db
}

// SetupDB - setup dabase for sharing to all api
func SetupDB() {

	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=disable password=%v TimeZone=%v", dbHost, port, username, dbName, password, TimeZone)

	database, err :=  gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}
	// defer database.Close()


	database.AutoMigrate(&model.User{})
	database.AutoMigrate(&model.Party{})

	db = database
}
