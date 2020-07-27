package mysql

import (
	"time"

	"github.com/arunvm/focus/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) GetOrganisationMembers(organisationID string) (*[]models.OrganisationMemberInfo, error) {
	var members []models.OrganisationMemberInfo

	err := db.Client.Table("organisation_members").Joins("JOIN users on organisation_members.user_id = users.id").
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

func (db *Mysql) CheckIfOrganisationMember(organisationID string, user *models.User) bool {
	return db.checkIfUserIsOrganisationMember(user.ID, organisationID)
}

func (db *Mysql) checkIfUserIsOrganisationMember(userID int, organisationID string) bool {
	var count int

	err := db.Client.Table("organisations").Joins("JOIN organisation_members on organisations.id = organisation_members.organisation_id").
		Where("organisations.id = ? AND archived = false AND organisation_members.user_id = ? AND organisations.type = ?", organisationID, userID, models.Organistation).
		Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":           "checkIfUserIsOrganisationMember",
			"info":           "checking if user is a member of organisation",
			"userID":         userID,
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
	orgMember := models.OrganisationMember{
		OrganisationID: orgID,
		UserID:         userID,
		JoinedAt:       time.Now().Unix(),
	}

	err := db.Create(orgMember).Error
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
