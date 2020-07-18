package mysql

import (
	"github.com/arunvm/travail-backend/models"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateOrganisationInviteToken(args *models.InviteToOrganisationArgs, admin *models.User) (*models.OrganisationInvitationInfo, error) {
	org := models.OrganisationInvitation{
		Email:          args.Email,
		OrganisationID: args.OrganisationID,
		Token:          xid.New().String(),
	}

	err := db.Client.Save(org).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateOrganisationInviteToken",
			"subFunc": "org.Save",
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	orgName, err := db.GetOrganisationName(org.OrganisationID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateOrganisationInviteToken",
			"subFunc": "db.GetOrganisationName",
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	return &models.OrganisationInvitationInfo{
		OrganisationID:   org.OrganisationID,
		Token:            org.Token,
		Email:            org.Email,
		OrganisationName: orgName,
	}, nil
}

func (db *Mysql) AcceptOrganisationInvite(args *models.AcceptOrganisationInviteArgs, user *models.User) error {
	var invite models.OrganisationInvitation

	tx := db.Client.Begin()
	err := tx.Find(&invite, "email = ? AND token = ?", user.Email, args.Token).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":   "AcceptOrganisationInvite",
			"info":   "checking if invitation to organisation exists",
			"userID": user.ID,
		}).Error(err)
		return err
	}

	err = addUserToOrganisation(tx, user.ID, invite.OrganisationID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":           "AcceptOrganisationInvite",
			"subFunc":        "addUserToOrganisation",
			"userID":         user.ID,
			"organisationID": invite.OrganisationID,
		}).Error(err)
		return err
	}

	err = tx.Delete(&invite).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":           "AcceptOrganisationInvite",
			"info":           "deleting invite",
			"userID":         user.ID,
			"organisationID": invite.OrganisationID,
		}).Error(err)
		return err
	}

	tx.Commit()
	return nil
}
