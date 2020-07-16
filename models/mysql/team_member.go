package mysql

import (
	"errors"
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) AddTeamMember(args models.AddTeamMemberArgs) error {
	team, err := db.getTeamFromID(args.TeamID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "AddTeamMember",
			"subFunc": "getTeamFromID",
			"teamID":  args.TeamID,
		}).Error(err)
		return err
	}

	if db.checkIfUserIsOrganisationMember(args.UserID, team.OrganisationID) == false {
		return errors.New("User not part of organisation")
	}

	member := models.TeamMember{
		TeamID:   team.ID,
		JoinedAt: time.Now().Unix(),
		UserID:   args.UserID,
	}

	err = db.Client.Create(member).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "AddTeamMember",
			"subFunc": "member.Create",
			"teamID":  args.TeamID,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) GetTeamMembers(teamID string) (*[]models.TeamMemberInfo, error) {
	var members []models.TeamMemberInfo

	err := db.Client.Table("team_members").Joins("JOIN users on team_members.user_id = users.id").
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
	teamMember := models.TeamMember{
		TeamID:   teamID,
		JoinedAt: time.Now().Unix(),
		UserID:   userID,
	}

	err := db.Create(teamMember).Error
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

func (db *Mysql) CheckIfTeamMember(teamID string, user *models.User) bool {
	var count int

	err := db.Client.Table("teams").Joins("JOIN team_members on teams.id = team_members.team_id").
		Joins("JOIN organisations on teams.organisation_id = organisations.id").
		Where("teams.id = ? AND teams.archived = false AND team_members.user_id = ? AND organisations.type = ?", teamID, user.ID, models.Organistation).
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
