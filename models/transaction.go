package models

type Transaction interface {
	Begin() Transaction
	Commit()
	Rollback()
	// Functions That require transactions
	UserSignup(args *SignUpArgs, googleOauth bool) (*User, string, error)
	CreateForgotPasswordToken(user *User) (string, error)
	CreateOrganisationInviteToken(args *InviteToOrganisationArgs, admin *User) (*OrganisationInvitationInfo, error)
	CheckIfUserExists(email string) bool
	GetUserFromEmail(email string) (*User, error)
	GetPersonalTeamID(user *User) (string, error)
	CreateEmailValidationToken(user *User) (string, error)
}
