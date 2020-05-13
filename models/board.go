package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Board struct {
	ID        string `json:"id" gorm:"primary_key;auto_increment:false"`
	AdminID   int    `json:"adminID"`
	TeamID    string `json:"teamID"`
	Title     string `json:"title"`
	CreatedOn int64  `json:"createdOn"`
	Archived  bool   `json:"archived"`
}

func (b *Board) Create(db *gorm.DB) error {
	return db.Create(&b).Error
}

func (b *Board) Save(db *gorm.DB) error {
	return db.Save(&b).Error
}

type CreateBoardArgs struct {
	TeamID string `json:"teamID"`
	Title  string `json:"title" binding:"required"`
}

type UpdateBoardArgs struct {
	ID     string  `json:"id" binding:"required"`
	TeamID string  `json:"teamID"`
	Title  *string `json:"title"`
}

func (teamMember *User) CreateBoard(db *gorm.DB, args *CreateBoardArgs) error {
	board := Board{
		ID:        uuid.New().String(),
		AdminID:   teamMember.ID,
		Archived:  false,
		CreatedOn: time.Now().Unix(),
		TeamID:    args.TeamID,
		Title:     args.Title,
	}

	err := board.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "CreateBoard",
			"subFunc":      "board.Create",
			"args":         *args,
			"teamMemberID": teamMember.ID,
		}).Error(err)
		return err
	}

	return nil
}

func (teamMember *User) UpdateBoard(db *gorm.DB, args *UpdateBoardArgs) error {
	board, err := getBoard(db, args.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "UpdateBoard",
			"subFunc":      "getBoard",
			"teamMemberID": teamMember.ID,
			"args":         *args,
		}).Error(err)
		return err
	}

	if board.TeamID != args.TeamID {
		return errors.New("User cannot access this board")
	}

	if args.Title != nil {
		board.Title = *args.Title
	}

	err = board.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "UpdateBoard",
			"subFunc":      "board.Save",
			"teamMemberID": teamMember.ID,
			"args":         *args,
		}).Error(err)
		return err
	}

	return nil
}

func getBoard(db *gorm.DB, boardID string) (*Board, error) {
	var board Board

	err := db.Find(&board, "id = ? AND archived = false", boardID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "getBoard",
			"info":    "retrieving board with id",
			"boardID": boardID,
		}).Error(err)
		return nil, err
	}

	return &board, nil
}
