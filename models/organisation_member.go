package models

type OrganisationMember struct {
	OrganisationID string `json:"organisationID" gorm:"primary_key;auto_increment:false"`
	UserID         int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	JoinedAt       int64  `json:"joinedAt"`
}

type OrganisationMemberInfo struct {
	OrganisationID string  `json:"-"`
	UserID         int     `json:"userId"`
	Name           string  `json:"name"`
	ProfilePic     *string `json:"profilePic"`
}
