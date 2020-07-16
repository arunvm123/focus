package mysql

import (
	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Mysql struct {
	Client *gorm.DB
}

func New(connectionString string) (*Mysql, error) {
	client, err := gorm.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	client.LogMode(true)

	return &Mysql{
		Client: client,
	}, nil
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&models.EmailValidateToken{})
	db.AutoMigrate(&models.ForgotPasswordToken{})
	db.AutoMigrate(&models.FCMNotificationToken{})

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.List{})
	db.AutoMigrate(&models.Task{})

	db.AutoMigrate(&models.Organisation{})
	db.AutoMigrate(&models.OrganisationMember{})
	db.AutoMigrate(&models.OrganisationInvitation{})

	db.AutoMigrate(&models.Team{})
	db.AutoMigrate(&models.TeamMember{})

	db.AutoMigrate(&models.Board{})
	db.AutoMigrate(&models.BoardColumn{})
	db.AutoMigrate(&models.ColumnCard{})

	db.AutoMigrate(models.Bug{})

	err := db.Model(models.List{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for list model\n%v", err)
	}
	err = db.Model(models.List{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for list model\n%v", err)
	}
	err = db.Model(models.Task{}).AddForeignKey("list_id", "lists(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for task model\n%v", err)
	}
	err = db.Model(models.EmailValidateToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for email_validate_token model\n%v", err)
	}
	err = db.Model(models.FCMNotificationToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for fcm_notification_token  model\n%v", err)
	}
	err = db.Model(models.ForgotPasswordToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for forgot_password_token model\n%v", err)
	}
	err = db.Model(models.Bug{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for bug model\n%v", err)
	}
	err = db.Model(models.Organisation{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisations model\n%v", err)
	}
	err = db.Model(models.OrganisationMember{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation members model\n%v", err)
	}
	err = db.Model(models.OrganisationMember{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation members model\n%v", err)
	}
	err = db.Model(models.Team{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
	err = db.Model(models.Team{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for teams model\n%v", err)
	}
	err = db.Model(models.TeamMember{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for team members model\n%v", err)
	}
	err = db.Model(models.TeamMember{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for team members model\n%v", err)
	}
	err = db.Model(models.OrganisationInvitation{}).AddForeignKey("organisation_id", "organisations(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for organisation_invitations model\n%v", err)
	}
	err = db.Model(models.Board{}).AddForeignKey("admin_id", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board model\n%v", err)
	}
	err = db.Model(models.Board{}).AddForeignKey("team_id", "teams(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board model\n%v", err)
	}
	err = db.Model(models.BoardColumn{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for board_column model\n%v", err)
	}
	err = db.Model(models.ColumnCard{}).AddForeignKey("column_id", "board_columns(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
	err = db.Model(models.ColumnCard{}).AddForeignKey("assigned_to", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
	err = db.Model(models.ColumnCard{}).AddForeignKey("assigned_by", "users(id)", "RESTRICT", "RESTRICT").Error
	if err != nil {
		log.Fatalf("Error adding foreign key for column_card model\n%v", err)
	}
}
