package models

type OrganisationInvitation struct {
	OrganisationID string `json:"organisationID" gorm:"primary_key"`
	Email          string `json:"email" gorm:"primary_key"`
	Token          string `json:"token"`
}

type InviteToOrganisationArgs struct {
	OrganisationID string `json:"organisationID" binding:"-"`
	Email          string `json:"email" binding:"required,email"`
}

type AcceptOrganisationInviteArgs struct {
	Token string `json:"token" binding:"required"`
}

type OrganisationInvitationInfo struct {
	OrganisationID   string `json:"organisationID"`
	Email            string `json:"email"`
	Token            string `json:"token"`
	OrganisationName string `json:"organisationName"`
}
