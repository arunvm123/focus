package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type OrganisationInvitation struct {
	OrganisationID string `json:"organisationID" gorm:"primary_key"`
	Email          string `json:"email" gorm:"primary_key"`
	Token          string `json:"token"`
}

// Create is a helper function to create a new organisation
func (oi *OrganisationInvitation) Create(db *gorm.DB) error {
	return db.Create(&oi).Error
}

// Save is a helper function to update existing organisation
func (oi *OrganisationInvitation) Save(db *gorm.DB) error {
	return db.Save(&oi).Error
}

type InviteToOrganisationArgs struct {
	OrganisationID string `json:"organisationID" binding:"-"`
	Email          string `json:"email" binding:"required,email"`
}

type AcceptOrganisationInviteArgs struct {
	Token string `json:"token" binding:"required"`
}

func (admin *User) CreateOrganisationInviteToken(db *gorm.DB, args *InviteToOrganisationArgs) (*OrganisationInvitation, error) {
	org := OrganisationInvitation{
		Email:          args.Email,
		OrganisationID: args.OrganisationID,
		Token:          xid.New().String(),
	}

	err := org.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateOrganisationInviteToken",
			"subFunc": "org.Save",
			"adminID": admin.ID,
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	return &org, nil
}

func (user *User) AcceptOrganisationInvite(db *gorm.DB, args *AcceptOrganisationInviteArgs) error {
	var invite OrganisationInvitation

	err := db.Find(&invite, "email = ? AND token = ?", user.Email, args.Token).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "AcceptOrganisationInvite",
			"info":   "checking if invitation to organisation exists",
			"userID": user.ID,
		}).Error(err)
		return err
	}

	err = addUserToOrganisation(db, user.ID, invite.OrganisationID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "AcceptOrganisationInvite",
			"subFunc":        "addUserToOrganisation",
			"userID":         user.ID,
			"organisationID": invite.OrganisationID,
		}).Error(err)
		return err
	}

	err = db.Delete(&invite).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "AcceptOrganisationInvite",
			"info":           "deleting invite",
			"userID":         user.ID,
			"organisationID": invite.OrganisationID,
		}).Error(err)
		return err
	}

	return nil
}
