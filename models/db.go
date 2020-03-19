package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&EmailValidateToken{})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&List{})
	db.AutoMigrate(&Task{})

	err := db.Model(List{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for list model\n%v", err)
	}
	err = db.Model(Task{}).AddForeignKey("list_id", "lists(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for task model\n%v", err)
	}
	err = db.Model(EmailValidateToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for email_validate_token  model\n%v", err)
	}
}
