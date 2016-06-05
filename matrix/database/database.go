package database

import (
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func DatabaseInit() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
		panic("Failed to create database")
	}

	db.AutoMigrate(&AccessToken{})
	db.AutoMigrate(&Account{})

	DB = db

	return db
}
