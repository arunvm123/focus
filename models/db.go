package models

import (
	push "github.com/arunvm/travail-backend/push_notification"
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
	CreateEmailValidationToken(user *User) (string, error)
	VerifyEmail(token string) error
	InvalidateEmailTokens(userID int) error
	// FCM notification token functions
	AddNotificationToken(args *AddNotificationTokenArgs, user *User) error
	GetNotificationTokens(user *User) ([]string, error)
	// Forgot password token functions
	CreateForgotPasswordToken(user *User) (string, error)
	ResetPassword(token, password string) error
	// List functions
	CreateList(args *CreateListArgs, user *User) (*List, error)
	GetLists(args *GetListsArgs, user *User) (*[]ListInfo, error)
	UpdateList(args *UpdateListArgs, user *User) error
	// Organisation invitation functions
	CreateOrganisationInviteToken(args *InviteToOrganisationArgs, admin *User) (*OrganisationInvitationInfo, error)
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
	UserSignup(args *SignUpArgs, googleOauth bool) (*User, string, error)
	UpdatePassword(args *UpdatePasswordArgs, user *User) error
	// Transaction
	Transaction
}
