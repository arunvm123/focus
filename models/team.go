package models

import "github.com/jinzhu/gorm"

type Team struct {
	ID             string  `json:"id" gorm:"primary_key;auto_increment:false"`
	OrganisationID string  `json:"organisationID"`
	AdminID        int     `json:"adminID"`
	Name           string  `json:"name"`
	Description    *string `json:"description" gorm:"size:3000"`
	CreatedAt      int64   `json:"createdAt"`
	Archived       bool    `json:"archived"`
}

// Create is a helper function to create a new organisation
func (team *Team) Create(db *gorm.DB) error {
	return db.Create(&team).Error
}

// Save is a helper function to update existing organisation
func (team *Team) Save(db *gorm.DB) error {
	return db.Save(&team).Error
}
