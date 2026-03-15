package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Email        string         `gorm:"not null;uniqueIndex"`
	PasswordHash string         `gorm:"column:password_hash;not null"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
