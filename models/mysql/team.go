package mysql

import (
	"time"

	"github.com/arunvm/focus/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateTeam(args *models.CreateTeamArgs, user *models.User) error {
	tx := db.Client.Begin()

	team := models.Team{
		ID:             uuid.New().String(),
		AdminID:        user.ID,
		Archived:       false,
		CreatedAt:      time.Now().Unix(),
		Name:           args.Name,
		OrganisationID: args.OrganisationID,
		Description:    args.Description,
	}

	err := tx.Create(team).Error
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateTeam",
			"subFunc": "team.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	err = addUserToTeam(tx, user.ID, team.ID)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateTeam",
			"subFunc": "addUserToTeam",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	tx.Commit()
	return nil
}

func (db *Mysql) UpdateTeam(args *models.UpdateTeamArgs, teamAdmin *models.User) error {
	team, err := db.getTeamFromID(args.TeamID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "UpdateTeam",
			"subFunc":     "getTeamFromID",
			"teamAdminID": teamAdmin.ID,
			"teamID":      args.TeamID,
		}).Error(err)
		return err
	}

	if args.Description != nil {
		team.Description = args.Description
	}
	if args.Name != nil {
		team.Name = *args.Name
	}

	err = db.Client.Save(team).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "UpdateTeam",
			"subFunc":     "team.Save",
			"teamAdminID": teamAdmin.ID,
			"teamID":      args.TeamID,
		}).Error(err)
		return err
	}

	return nil
}

func createPersonalTeam(db *gorm.DB, org *models.Organisation, user *models.User) error {
	team := models.Team{
		ID:             uuid.New().String(),
		AdminID:        user.ID,
		Archived:       false,
		CreatedAt:      time.Now().Unix(),
		Name:           models.PersonalString,
		OrganisationID: org.ID,
	}

	err := db.Create(team).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createPersonalTeam",
			"subFunc": "team.Create",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	err = addUserToTeam(db, user.ID, team.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createPersonalTeam",
			"subFunc": "addUserToTeam",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) GetPersonalTeamID(user *models.User) (string, error) {
	var team models.Team

	err := db.Client.Table("organisations").Joins("JOIN teams on organisations.id = teams.organisation_id").
		Where("organisations.admin_id = ? AND type = ?", user.ID, models.Personal).Select("teams.*").
		Find(&team).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "GetPersonalTeamID",
			"info": "retreiving personal team of user",
		}).Error(err)
		return "", err
	}

	return team.ID, nil
}

func (db *Mysql) getTeamFromID(teamID string) (*models.Team, error) {
	var team models.Team

	err := db.Client.Find(&team, "id = ? AND archived = false", teamID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "getTeamFromID",
			"info":   "retrieving team with ID",
			"teamID": teamID,
		}).Error(err)
		return nil, err
	}

	return &team, nil
}

func (db *Mysql) CheckIfTeamAdmin(teamID string, user *models.User) bool {
	var count int

	err := db.Client.Table("teams").Joins("JOIN organisations on teams.organisation_id = organisations.id").
		Where("teams.id = ? AND teams.archived = false AND teams.admin_id = ? AND organisations.type = ?", teamID, user.ID, models.Organistation).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "checkIfTeamAdmin",
			"info":   "checking if user id team admin",
			"userID": user.ID,
			"teamID": teamID,
		}).Error(err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
}
