package mysql

import (
	"errors"

	"github.com/arunvm/focus/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateBoardColumn(args *models.CreateBoardColumnArgs) error {
	var bc models.BoardColumn

	bc.ID = uuid.New().String()
	bc.BoardID = args.BoardID
	bc.Name = args.Name

	err := db.Client.Create(&bc).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "bc.Create",
			"subFunc": "bc.Create",
			"boardID": args.BoardID,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) GetBoardColumns(boardID string) (*[]models.BoardColumn, error) {
	var bc []models.BoardColumn

	err := db.Client.Find(&bc, "board_id = ?", boardID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "GetBoardColumns",
			"info":    "retrieving columns of specified board",
			"boardID": boardID,
		}).Error(err)
		return nil, err
	}

	return &bc, nil
}

func (db *Mysql) UpdateBoardColumn(args *models.UpdateBoardColumnArgs) error {
	column, err := db.getBoardColumn(args.ColumnID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "UpdateBoardColumn",
			"subFunc":  "getBoardColumn",
			"columnID": args.ColumnID,
		}).Error(err)
		return err
	}

	if column.BoardID != args.BoardID {
		return errors.New("Column not part of specified board")
	}

	if args.Name != nil {
		column.Name = *args.Name
	}

	err = db.Client.Save(&column).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "UpdateBoardColumn",
			"subFunc":  "column.Save",
			"columnID": args.ColumnID,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) CheckIfColumnPartOfBoard(boardColumnID string, boardID string) bool {
	boardColumn, err := db.getBoardColumn(boardColumnID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":         "CheckIfColumnPartOfBoard",
			"subFunc":      "getBoardColumn",
			"boardColunID": boardColumnID,
		}).Error(err)
		return false
	}

	if boardColumn.BoardID != boardID {
		return false
	}

	return true
}

func (db *Mysql) getBoardColumn(columnID string) (*models.BoardColumn, error) {
	var column models.BoardColumn

	err := db.Client.Find(&column, "id = ?", columnID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "getBoardColumn",
			"info":     "retrieving column details",
			"columdID": columnID,
		}).Error(err)
		return nil, err
	}

	return &column, nil
}
