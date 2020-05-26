package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type ColumnCard struct {
	ID          string `json:"id" gorm:"primary_key"`
	ColumnID    string `json:"columnID"`
	Heading     string `json:"heading"`
	Description string `json:"description" gorm:"size:3000"`
	AssignedTo  *int   `json:"assignedTo"`
	AssignedBy  *int   `json:"AssignedBy"`
	AssignedOn  *int64 `json:"assignedOn"`
}

func (cc *ColumnCard) Create(db *gorm.DB) error {
	return db.Create(&cc).Error
}

func (cc *ColumnCard) Save(db *gorm.DB) error {
	return db.Save(&cc).Error
}

type CreateColumnCardArgs struct {
	ColumnID    string `json:"columnID"`
	Heading     string `json:"heading" binding:"required"`
	Description string `json:"description"`
	AssignedTo  *int   `json:"assignedTo"`
	AssignedBy  *int   `json:"AssignedBy"`
}

func CreateColumnCard(db *gorm.DB, args *CreateColumnCardArgs) error {
	cc := ColumnCard{
		ColumnID:    args.ColumnID,
		AssignedBy:  args.AssignedBy,
		AssignedTo:  args.AssignedTo,
		Description: args.Description,
		Heading:     args.Heading,
		ID:          uuid.New().String(),
	}

	if cc.AssignedTo != nil {
		now := time.Now().Unix()
		cc.AssignedOn = &now
	}

	err := cc.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "CreateColumnCard",
			"subFunc":  "cc.Create",
			"columnID": args.ColumnID,
		}).Error(err)
		return err
	}

	return nil
}
