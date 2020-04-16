package models

import "github.com/jinzhu/gorm"

// Organisation groups all the teams together and is typically the company name
type Organisation struct {
	ID        string `json:"id" gorm:"primary_key;auto_increment:false"`
	AdminID   int    `json:"adminID"`
	Name      string `json:"name"`
	Type      int    `json:"type"` // Type denotes if this is the user's personal space or of a companies
	CreatedAt int64  `json:"createdAt"`
	Archived  bool   `json:"archived"`
}

const (
	PERSONAL     = 1
	ORGANISATION = 2
)

// Create is a helper function to create a new organisation
func (org *Organisation) Create(db *gorm.DB) error {
	return db.Create(&org).Error
}

// Save is a helper function to update existing organisation
func (org *Organisation) Save(db *gorm.DB) error {
	return db.Save(&org).Error
}
