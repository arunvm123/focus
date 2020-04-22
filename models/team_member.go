package models

import (
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
