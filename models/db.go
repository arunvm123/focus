package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&EmailValidateToken{})
	db.AutoMigrate(&ForgotPasswordToken{})
	db.AutoMigrate(&FCMNotificationToken{})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&List{})
	db.AutoMigrate(&Task{})

	db.AutoMigrate(&Organisation{})
	db.AutoMigrate(&Team{})

	db.AutoMigrate(Bug{})

	err := db.Model(List{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for list model\n%v", err)
	}
	err = db.Model(List{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for list model\n%v", err)
	}
	err = db.Model(Task{}).AddForeignKey("list_id", "lists(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for task model\n%v", err)
	}
	err = db.Model(EmailValidateToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for email_validate_token model\n%v", err)
	}
	err = db.Model(FCMNotificationToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for fcm_notification_token  model\n%v", err)
	}
	err = db.Model(ForgotPasswordToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for forgot_password_token model\n%v", err)
	}
	err = db.Model(Bug{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for bug model\n%v", err)
	}
	err = db.Model(Organisation{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisations model\n%v", err)
	}
	err = db.Model(Team{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
	err = db.Model(Team{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
}
