package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

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
	personal       = 1
	organistation  = 2
	personalString = "Personal"
)

// Create is a helper function to create a new organisation
func (org *Organisation) Create(db *gorm.DB) error {
	return db.Create(&org).Error
}

// Save is a helper function to update existing organisation
func (org *Organisation) Save(db *gorm.DB) error {
	return db.Save(&org).Error
}

type CreateOrganisationArgs struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOrganisationArgs struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

func (user *User) CreateOrganisation(db *gorm.DB, args *CreateOrganisationArgs) error {
	org := Organisation{
		ID:        uuid.NewV4().String(),
		AdminID:   user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Name:      args.Name,
		Type:      organistation,
	}

	err := org.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateOrganisation",
			"subFunc": "org.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	err = addUserToOrganisation(db, user.ID, org.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "CreateOrganisation",
			"subFunc":        "addUserToOrganisation",
			"userID":         user.ID,
			"organisationID": org.ID,
			"args":           *args,
		}).Error(err)
		return err
	}

	return nil
}

func (admin *User) UpdateOrganisation(db *gorm.DB, args *UpdateOrganisationArgs) error {
	org, err := getOrganisationFromID(db, args.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "UpdateOrganisation",
			"subFunc":        "getOrganisationFromID",
			"adminID":        admin.ID,
			"organisationID": args.ID,
		}).Error(err)
		return err
	}

	org.Name = args.Name
	err = org.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "UpdateOrganisation",
			"subFunc":        "org.Save",
			"adminID":        admin.ID,
			"organisationID": org.ID,
		}).Error(err)
		return err
	}

	return nil
}

func getOrganisationFromID(db *gorm.DB, organisationID string) (*Organisation, error) {
	var org Organisation

	err := db.Find(&org, "id = ? AND archived = false", organisationID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "getOrganisationFromID",
			"info":           "retrieving organisation specified by id",
			"organisationID": organisationID,
		}).Error(err)
		return nil, err
	}

	return &org, nil
}

func (user *User) createPersonalOrganisation(db *gorm.DB) (*Organisation, error) {
	org := Organisation{
		ID:        uuid.NewV4().String(),
		AdminID:   user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Name:      personalString,
		Type:      personal,
	}

	err := org.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createPersonalOrganisation",
			"subFunc": "org.Create",
			"userID":  user.ID,
		}).Error(err)
		return nil, err
	}

	return &org, nil
}

// In the case of an error, returns false and assumes that the user is not admin
func (user *User) CheckIfOrganisationAdmin(db *gorm.DB, orgID string) bool {
	var count int

	err := db.Table("organisations").Where("id = ? AND admin_id = ? AND archived = false AND type = ?", orgID, user.ID, organistation).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "CheckIfOrganisationAdmin",
			"info":           "checking if organisation exists with given user as admin",
			"userID":         user.ID,
			"organisationID": orgID,
		}).Error(err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
}
