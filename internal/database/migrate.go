package database

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string    `json:"name"     gorm:"not null"`
	Email    string    `json:"email"    gorm:"not null;unique"`
	Password string    `json:"password" gorm:"not null"`
}

func Migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(&User{})
}
