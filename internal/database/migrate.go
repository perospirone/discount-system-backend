package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name          string    `json:"name"     gorm:"not null"`
	Email         string    `json:"email"    gorm:"not null;unique"`
	Password      string    `json:"password" gorm:"not null"`
	HasDiscount   bool      `json:"has_discount" gorm:"not null"` // todo: put false as default
	AffiliateUserID uint      `json:"affiliate_user_id"`
	DiscountUse   time.Time `json:"discount_use"`
}

func Migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(&User{})
}
