package models

// Organisation groups all the teams together and is typically the company name
type Organisation struct {
	ID             string  `json:"id" gorm:"primary_key;auto_increment:false"`
	AdminID        int     `json:"adminID"`
	Name           string  `json:"name"`
	DisplayPicture *string `json:"displayPicture"`
	Theme          string  `json:"theme"`
	Type           int     `json:"type"` // Type denotes if this is the user's personal space or of a companies
	CreatedAt      int64   `json:"createdAt"`
	Archived       bool    `json:"archived"`
}

const (
	Personal       = 1
	Organistation  = 2
	PersonalString = "Personal"
)

type CreateOrganisationArgs struct {
	Name           string  `json:"name" binding:"required"`
	DisplayPicture *string `json:"displayPicture"`
	Theme          string  `json:"theme" binding:"required"`
}

type UpdateOrganisationArgs struct {
	ID             string  `json:"-"`
	Name           *string `json:"name"`
	DisplayPicture *string `json:"displayPicture"`
	Theme          *string `json:"theme"`
}
