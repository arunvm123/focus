package mysql

import (
	"errors"
	"time"

	"github.com/arunvm/focus/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateOrganisation(args *models.CreateOrganisationArgs, user *models.User) error {
	tx := db.Client.Begin()

	org := models.Organisation{
		ID:             uuid.New().String(),
		AdminID:        user.ID,
		Archived:       false,
		CreatedAt:      time.Now().Unix(),
		Name:           args.Name,
		Type:           models.Organistation,
		DisplayPicture: args.DisplayPicture,
		Theme:          args.Theme,
	}

	err := tx.Create(org).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateOrganisation",
			"subFunc": "org.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	err = addUserToOrganisation(tx, user.ID, org.ID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":           "CreateOrganisation",
			"subFunc":        "addUserToOrganisation",
			"userID":         user.ID,
			"organisationID": org.ID,
			"args":           *args,
		}).Error(err)
		return err
	}

	tx.Commit()
	return nil
}

func (db *Mysql) GetOrganisations(user *models.User) (*[]models.Organisation, error) {
	var organisations []models.Organisation

	err := db.Client.Table("organisations").Joins("JOIN organisation_members on organisations.id = organisation_members.organisation_id").
		Where("user_id = ? AND organisations.archived = false", user.ID).Select("organisations.*").Find(&organisations).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetOrganisations",
			"info":   "retrieving organisations of user",
			"userID": user.ID,
		}).Error(err)
		return nil, err
	}

	return &organisations, nil
}

func (db *Mysql) UpdateOrganisation(args *models.UpdateOrganisationArgs, admin *models.User) error {
	org, err := db.getOrganisationFromID(args.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "UpdateOrganisation",
			"subFunc":        "getOrganisationFromID",
			"adminID":        admin.ID,
			"organisationID": args.ID,
		}).Error(err)
		return err
	}

	if args.Name != nil {
		org.Name = *args.Name
	}
	if args.DisplayPicture != nil {
		org.DisplayPicture = args.DisplayPicture
	}
	if args.Theme != nil {
		org.Theme = *args.Theme
	}

	err = db.Client.Save(org).Error
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

func (db *Mysql) getOrganisationFromID(organisationID string) (*models.Organisation, error) {
	var org models.Organisation

	err := db.Client.Find(&org, "id = ? AND archived = false", organisationID).Error
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

func createPersonalOrganisation(db *gorm.DB, user *models.User) (*models.Organisation, error) {
	org := models.Organisation{
		ID:        uuid.New().String(),
		AdminID:   user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Name:      models.PersonalString,
		Type:      models.Personal,
	}

	err := db.Create(org).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createPersonalOrganisation",
			"subFunc": "org.Create",
			"userID":  user.ID,
		}).Error(err)
		return nil, err
	}

	err = addUserToOrganisation(db, user.ID, org.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UserSignup",
			"subFunc": "addUserToOrganisation",
			"userID":  user.ID,
		}).Error(err)
		return nil, err
	}

	return &org, nil
}

// In the case of an error, returns false and assumes that the user is not admin
func (db *Mysql) CheckIfOrganisationAdmin(orgID string, user *models.User) bool {
	var count int

	err := db.Client.Table("organisations").Where("id = ? AND admin_id = ? AND archived = false AND type = ?", orgID, user.ID, models.Organistation).Count(&count).Error
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

func (db *Mysql) GetOrganisationName(organisationID string) (string, error) {
	var name []string

	err := db.Client.Table("organisations").Where("id = ? AND archived = false", organisationID).Pluck("name", &name).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "GetOrganisationName",
			"info":           "retrieving organisation name",
			"organisationID": organisationID,
		}).Error(err)
		return "", err
	}

	if len(name) != 1 {
		return "", errors.New("error fetching name of organiation")
	}

	return name[0], nil
}
