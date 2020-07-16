package mysql

import (
	"github.com/arunvm/travail-backend/email"
	"github.com/arunvm/travail-backend/models"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateOrganisationInviteToken(args *models.InviteToOrganisationArgs, admin *models.User, emailClient email.Email) error {
	org := models.OrganisationInvitation{
		Email:          args.Email,
		OrganisationID: args.OrganisationID,
		Token:          xid.New().String(),
	}

	tx := db.Client.Begin()
	err := tx.Save(org).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateOrganisationInviteToken",
			"subFunc": "org.Save",
			"args":    *args,
		}).Error(err)
		return err
	}

	orgName, err := db.GetOrganisationName(org.OrganisationID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateOrganisationInviteToken",
			"subFunc": "db.GetOrganisationName",
			"args":    *args,
		}).Error(err)
		return err
	}

	err = emailClient.SendOrganisationInvite(admin.Name, org.Email, org.Token, orgName)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "inviteToOrganisation",
			"subFunc": "emails.SendInviteToOrganisation",
			"adminID": admin.ID,
			"args":    args,
		}).Error(err)
		return err
	}

	tx.Commit()
	return nil
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
