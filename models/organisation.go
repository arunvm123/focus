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
