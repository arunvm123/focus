package models

import (
	"github.com/arunvm/travail-backend/email"
	push "github.com/arunvm/travail-backend/push_notification"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type DB interface {
	// Board Column functions
	CreateBoardColumn(args *CreateBoardColumnArgs) error
	GetBoardColumns(boardID string) (*[]BoardColumn, error)
	UpdateBoardColumn(args *UpdateBoardColumnArgs) error
	CheckIfColumnPartOfBoard(boardColumnID string, boardID string) bool
	// Board functions
	CreateBoard(args *CreateBoardArgs, teamMember *User) error
	UpdateBoard(args *UpdateBoardArgs) error
	GetBoards(args *GetBoardsArgs) (*[]Board, error)
	CheckIfBoardPartOfTeam(boardID, teamID string) bool
	// Bug functions
	CreateBug(args *CreateBugArgs, user *User) error
	GetBugs() (*[]BugInfo, error)
	UpdateBug(args *UpdateBugArgs) error
	// Column card functions
	CreateColumnCard(args *CreateColumnCardArgs) error
	GetColumnCards(columnID string) (*[]ColumnCard, error)
	UpdateColumnCard(args *UpdateColumnCardArgs, user *User) error
	// Email validate token functions
	CreateEmailValidationToken(user *User, emailCLient email.Email) error
	VerifyEmail(token string) error
	InvalidateEmailTokens(userID int) error
	// FCM notification token functions
	AddNotificationToken(args *AddNotificationTokenArgs, user *User) error
	GetNotificationTokens(user *User) ([]string, error)
	// Forgot password token functions
	CreateForgotPasswordToken(user *User, emailCLient email.Email) error
	ResetPassword(token, password string) error
	// List functions
	CreateList(args *CreateListArgs, user *User) (*List, error)
	GetLists(args *GetListsArgs, user *User) (*[]ListInfo, error)
	UpdateList(args *UpdateListArgs, user *User) error
	// Organisation invitation functions
	CreateOrganisationInviteToken(args *InviteToOrganisationArgs, admin *User, emailClient email.Email) error
	AcceptOrganisationInvite(args *AcceptOrganisationInviteArgs, user *User) error
	// Organisation member functions
	GetOrganisationMembers(organisationID string) (*[]OrganisationMemberInfo, error)
	CheckIfOrganisationMember(organisationID string, user *User) bool
	// Organisation functions
	CreateOrganisation(args *CreateOrganisationArgs, user *User) error
	GetOrganisations(user *User) (*[]Organisation, error)
	UpdateOrganisation(args *UpdateOrganisationArgs, admin *User) error
	CheckIfOrganisationAdmin(orgID string, user *User) bool
	GetOrganisationName(organisationID string) (string, error)
	// Task functions
	CreateTask(args *CreateTaskArgs, user *User) (*Task, error)
	GetTasks(args *GetTasksArgs, user *User) (*[]TaskInfo, error)
	UpdateTask(args UpdateTaskArgs, user *User) error
	DeleteTasks(args *DeleteTasksArgs, user *User) error
	// Team member functions
	AddTeamMember(args AddTeamMemberArgs) error
	GetTeamMembers(teamID string) (*[]TeamMemberInfo, error)
	CheckIfTeamMember(teamID string, user *User) bool
	// Team functions
	CreateTeam(args *CreateTeamArgs, user *User) error
	UpdateTeam(args *UpdateTeamArgs, teamAdmin *User) error
	GetPersonalTeamID(user *User) (string, error)
	CheckIfTeamAdmin(teamID string, user *User) bool
	SendPushNotificationForTasksAboutToExpire(pushClient push.Notification) error
	// User functions
	GetProfile(user *User) (*UserProfile, error)
	UpdateProfile(args UpdateProfileArgs, user *User) error
	GetUserFromEmail(email string) (*User, error)
	GetUserFromID(userID int) (*User, error)
	CheckIfUserExists(email string) bool
	UserSignup(args *SignUpArgs, googleOauth bool, emailClient email.Email) (*User, error)
	UpdatePassword(args *UpdatePasswordArgs, user *User) error
}

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
