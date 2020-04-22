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

type OrganisationMemberInfo struct {
	OrganisationID string  `json:"-"`
	UserID         int     `json:"userId"`
	Name           string  `json:"name"`
	ProfilePicture *string `json:"profilePicture"`
}

func GetOrganisationMembers(db *gorm.DB, organisationID string) (*[]OrganisationMemberInfo, error) {
	var members []OrganisationMemberInfo

	err := db.Table("organisation_members").Joins("JOIN users on organisation_members.user_id = users.id").
		Select("organisation_members.*,users.name,users.profile_pic").
		Where("organisation_members.organisation_id = ?", organisationID).
		Find(&members).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "GetOrganisationMembers",
			"info":           "retrieving info of organisation members",
			"organisationID": organisationID,
		}).Error(err)
		return nil, err
	}

	return &members, nil
}

func (user *User) CheckIfOrganisationMember(db *gorm.DB, organisationID string) bool {
	var count int

	err := db.Table("organisations").Joins("JOIN organisation_members on organisations.id = organisation_members.organisation_id").
		Where("organisations.id = ? AND archived = false AND organisation_members.user_id = ?", organisationID, user.ID).
		Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "CheckIfOrganisationMember",
			"info":           "checking if user is a member of organisation",
			"userID":         user.ID,
			"organisationID": organisationID,
		}).Error(err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
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
