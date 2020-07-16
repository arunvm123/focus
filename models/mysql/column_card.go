package mysql

import (
	"errors"
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateColumnCard(args *models.CreateColumnCardArgs) error {
	cc := models.ColumnCard{
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

	err := db.Client.Create(cc).Error
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

func (db *Mysql) GetColumnCards(columnID string) (*[]models.ColumnCard, error) {
	var cards []models.ColumnCard

	err := db.Client.Find(&cards, "column_id = ?", columnID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "GetColumnCards",
			"info":     "retrieving cards of the column",
			"columnID": columnID,
		}).Error(err)
		return nil, err
	}

	return &cards, nil
}

func (db *Mysql) UpdateColumnCard(args *models.UpdateColumnCardArgs, user *models.User) error {
	card, err := db.getColumnCard(args.CardID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "UpdateColumnCard",
			"subFunc":      "getColumnCard",
			"columnCardID": args.CardID,
		}).Error(err)
		return err
	}

	if card.ColumnID != args.ColumnID {
		return errors.New("Card does not beling to specified column")
	}

	if args.Description != nil {
		card.Description = *args.Description
	}
	if args.Heading != nil {
		card.Heading = *args.Heading
	}
	if args.AssignedTo != nil {
		card.AssignedTo = args.AssignedTo
		card.AssignedBy = &user.ID
		assignedOn := time.Now().Unix()
		card.AssignedOn = &assignedOn
	}

	err = db.Client.Save(card).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "UpdateColumnCard",
			"subFunc":      "card.Save",
			"userID":       user.ID,
			"columnCardID": card.ID,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) getColumnCard(cardID string) (*models.ColumnCard, error) {
	var card models.ColumnCard

	err := db.Client.Find(&card, "id = ?", cardID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "getColumnCard",
			"info":         "retrieving column card details",
			"columnCardID": cardID,
		}).Error(err)
		return nil, err
	}

	return &card, nil
}
