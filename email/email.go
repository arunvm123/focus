package email

type Email interface {
	SendValidationEmail(name, email string, token string) error
	SendForgotPasswordEmail(name, email string, token string) error
	SendOrganisationInvite(adminName, inviteEmail, token, orgName string) error
}
