package models

type Team struct {
	ID             string  `json:"id" gorm:"primary_key;auto_increment:false"`
	OrganisationID string  `json:"organisationID"`
	AdminID        int     `json:"adminID"`
	Name           string  `json:"name"`
	Description    *string `json:"description" gorm:"size:3000"`
	CreatedAt      int64   `json:"createdAt"`
	Archived       bool    `json:"archived"`
}

type CreateTeamArgs struct {
	OrganisationID string  `json:"-"`
	Name           string  `json:"name" binding:"required"`
	Description    *string `json:"description"`
}

type UpdateTeamArgs struct {
	TeamID      string  `json:"-"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
