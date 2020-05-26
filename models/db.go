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
	db.AutoMigrate(&OrganisationMember{})
	db.AutoMigrate(&OrganisationInvitation{})

	db.AutoMigrate(&Team{})
	db.AutoMigrate(&TeamMember{})

	db.AutoMigrate(&Board{})
	db.AutoMigrate(&BoardColumn{})
	db.AutoMigrate(&ColumnCard{})

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
	err = db.Model(OrganisationMember{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation members model\n%v", err)
	}
	err = db.Model(OrganisationMember{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation members model\n%v", err)
	}
	err = db.Model(Team{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
	err = db.Model(Team{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
	err = db.Model(TeamMember{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for team members model\n%v", err)
	}
	err = db.Model(TeamMember{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for team members model\n%v", err)
	}
	err = db.Model(OrganisationInvitation{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation_invitations model\n%v", err)
	}
	err = db.Model(Board{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board model\n%v", err)
	}
	err = db.Model(Board{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board model\n%v", err)
	}
	err = db.Model(BoardColumn{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board_column model\n%v", err)
	}
	err = db.Model(ColumnCard{}).AddForeignKey("column_id", "board_columns(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
	err = db.Model(ColumnCard{}).AddForeignKey("assigned_to", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
	err = db.Model(ColumnCard{}).AddForeignKey("assigned_by", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
}
