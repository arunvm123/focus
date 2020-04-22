package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type TeamMember struct {
	TeamID   string `json:"teamID" gorm:"primary_key;auto_increment:false"`
	UserID   int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	JoinedAt int64  `json:"joinedAt"`
}

func (tm *TeamMember) Create(db *gorm.DB) error {
	return db.Create(&tm).Error
}

func (tm *TeamMember) Save(db *gorm.DB) error {
	return db.Save(&tm).Error
}

type AddTeamMemberArgs struct {
	TeamID string `json:"-"`
	UserID int    `json:"userID" binding:"required"`
}

func (teamAdmin *User) AddTeamMember(db *gorm.DB, args AddTeamMemberArgs) error {
	team, err := getTeamFromID(db, args.TeamID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "AddTeamMember",
			"subFunc":     "getTeamFromID",
			"teamAdminID": teamAdmin.ID,
			"teamID":      args.TeamID,
		}).Error(err)
		return err
	}

	if checkIfUserIsOrganisationMember(db, args.UserID, team.OrganisationID) == false {
		return errors.New("User not part of organisation")
	}

	member := TeamMember{
		TeamID:   team.ID,
		JoinedAt: time.Now().Unix(),
		UserID:   args.UserID,
	}

	err = member.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "AddTeamMember",
			"subFunc":     "member.Create",
			"teamAdminID": teamAdmin.ID,
			"teamID":      args.TeamID,
		}).Error(err)
		return err
	}

	return nil
}

func addUserToTeam(db *gorm.DB, userID int, teamID string) error {
	teamMember := TeamMember{
		TeamID:   teamID,
		JoinedAt: time.Now().Unix(),
		UserID:   userID,
	}

	err := teamMember.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "addUserToTeam",
			"subFunc": "teamMember.Create",
			"userID":  userID,
			"teamID":  teamID,
		}).Error(err)
		return err
	}

	return nil
}
