package mysql

import (
	"errors"
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateBoard(args *models.CreateBoardArgs, teamMember *models.User) error {
	board := models.Board{
		ID:        uuid.New().String(),
		AdminID:   teamMember.ID,
		Archived:  false,
		CreatedOn: time.Now().Unix(),
		TeamID:    args.TeamID,
		Title:     args.Title,
	}

	err := db.Client.Create(&board).Error
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

func (db *Mysql) UpdateBoard(args *models.UpdateBoardArgs) error {
	board, err := db.getBoard(args.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateBoard",
			"subFunc": "getBoard",
			"args":    *args,
		}).Error(err)
		return err
	}

	if board.TeamID != args.TeamID {
		return errors.New("User cannot access this board")
	}

	if args.Title != nil {
		board.Title = *args.Title
	}

	err = db.Client.Save(&board).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateBoard",
			"subFunc": "board.Save",
			"args":    *args,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) GetBoards(args *models.GetBoardsArgs) (*[]models.Board, error) {
	var boards []models.Board

	err := db.Client.Find(&boards, "team_id = ? AND archived = false", args.TeamID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetBoards",
			"info":   "retrieving boards of team",
			"teamID": args.TeamID,
		}).Error(err)
		return nil, err
	}

	return &boards, nil
}

func (db *Mysql) CheckIfBoardPartOfTeam(boardID, teamID string) bool {
	board, err := db.getBoard(boardID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CheckIfBoardPartOfTeam",
			"subFunc": "getBoard",
			"boardID": boardID,
		}).Error(err)
		return false
	}

	if board.TeamID != teamID {
		return false
	}

	return true
}

func (db *Mysql) getBoard(boardID string) (*models.Board, error) {
	var board models.Board

	err := db.Client.Find(&board, "id = ? AND archived = false", boardID).Error
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
