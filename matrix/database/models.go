package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)


type Account struct {
	gorm.Model
	Username string
	Password string
}

type AccessToken struct {
	gorm.Model
	Token string
	Account Account
	CreatedAt time.Time
	ExpiresAt time.Time
}
