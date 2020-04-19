package models

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type OrganisationMember struct {
	OrganisationID string `json:"organisationID" gorm:"primary_key;auto_increment:false"`
	UserID         int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	JoinedAt       int64  `json:"joinedAt"`
}

// Create is a helper function to add a user to organisation
func (member *OrganisationMember) Create(db *gorm.DB) error {
	return db.Create(&member).Error
}

// Save is a helper function to update existing organisation member
func (member *OrganisationMember) Save(db *gorm.DB) error {
	return db.Save(&member).Error
}

func addUserToOrganisation(db *gorm.DB, userID int, orgID string) error {
	orgMember := OrganisationMember{
		OrganisationID: orgID,
		UserID:         userID,
		JoinedAt:       time.Now().Unix(),
	}

	err := orgMember.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "addUserToOrganisation",
			"subFunc": "orgMember.Create",
			"userID":  userID,
			"orgID":   orgID,
		}).Error(err)
		return err
	}

	return nil
}
