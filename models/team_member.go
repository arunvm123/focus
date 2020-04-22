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

type TeamMemberInfo struct {
	TeamID     string  `json:"teamID"`
	UserID     int     `json:"userID"`
	Name       string  `json:"name"`
	ProfilePic *string `json:"profile_pic"`
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

func GetTeamMembers(db *gorm.DB, teamID string) (*[]TeamMemberInfo, error) {
	var members []TeamMemberInfo

	err := db.Table("team_members").Joins("JOIN users on team_members.user_id = users.id").
		Select("team_members.*,users.name,users.profile_pic").
		Where("team_members.team_id = ?", teamID).
		Find(&members).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetTeamMembers",
			"info":   "retrieving info of team members",
			"teamID": teamID,
		}).Error(err)
		return nil, err
	}

	return &members, nil
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

func (user *User) CheckIfTeamMember(db *gorm.DB, teamID string) bool {
	var count int

	err := db.Table("teams").Joins("JOIN team_members on teams.id = team_members.team_id").
		Joins("JOIN organisations on teams.organisation_id = organisations.id").
		Where("teams.id = ? AND teams.archived = false AND team_members.user_id = ? AND organisations.type = ?", teamID, user.ID, organistation).
		Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "CheckIfTeamMember",
			"info":   "checking if user is a member of team",
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
