package db

import (
	"fmt"
	//"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	//dburl := os.Getenv("DATBASE_URL")
	var err error 
	dsn := "host=localhost user=postgres password=search dbname=postgres port=5432 sslmode=disable"
	DBConn, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		fmt.Println("failed to connect to db")
		panic("faied to connect to db")
	} else{
		fmt.Println("connection suceesful")
	}

	//enable uuid-ossp ext
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		panic(err)
	}
	
	err = DBConn.AutoMigrate(&User{}, &SearchSettings{})
	if err != nil {
		fmt.Println("Failed to migrate")
		panic(err)
	}
	
}
func GetDB() *gorm.DB {
	return DBConn
}