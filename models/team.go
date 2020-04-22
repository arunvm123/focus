package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Team struct {
	ID             string  `json:"id" gorm:"primary_key;auto_increment:false"`
	OrganisationID string  `json:"organisationID"`
	AdminID        int     `json:"adminID"`
	Name           string  `json:"name"`
	Description    *string `json:"description" gorm:"size:3000"`
	CreatedAt      int64   `json:"createdAt"`
	Archived       bool    `json:"archived"`
}

// Create is a helper function to create a new organisation
func (team *Team) Create(db *gorm.DB) error {
	return db.Create(&team).Error
}

// Save is a helper function to update existing organisation
func (team *Team) Save(db *gorm.DB) error {
	return db.Save(&team).Error
}

type CreateTeamArgs struct {
	OrganisationID string  `json:"-"`
	Name           string  `json:"name" binding:"required"`
	Description    *string `json:"description"`
}

func (user *User) CreateTeam(db *gorm.DB, args *CreateTeamArgs) error {
	team := Team{
		ID:             uuid.NewV4().String(),
		AdminID:        user.ID,
		Archived:       false,
		CreatedAt:      time.Now().Unix(),
		Name:           args.Name,
		OrganisationID: args.OrganisationID,
		Description:    args.Description,
	}

	err := team.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateTeam",
			"subFunc": "team.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	err = addUserToTeam(db, user.ID, team.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateTeam",
			"subFunc": "addUserToTeam",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	return nil
}

func (user *User) createPersonalTeam(db *gorm.DB, org *Organisation) error {
	team := Team{
		ID:             uuid.NewV4().String(),
		AdminID:        user.ID,
		Archived:       false,
		CreatedAt:      time.Now().Unix(),
		Name:           personalString,
		OrganisationID: org.ID,
	}

	err := team.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "createPersonalTeam",
			"subFunc": "team.Create",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}

func getTeamFromID(db *gorm.DB, teamID string) (*Team, error) {
	var team Team

	err := db.Find(&team, "id = ? AND archived = false", teamID).Error
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
