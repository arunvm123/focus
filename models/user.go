package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique;not null"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create(db *gorm.DB) error {
	return db.Create(&user).Error
}
