package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type BoardColumn struct {
	ID      string `json:"id" gorm:"primary_key;auto_increment:false"`
	BoardID string `json:"boardID"`
	Name    string `json:"name"`
}

func (bc *BoardColumn) Create(db *gorm.DB) error {
	return db.Create(&bc).Error
}

func (bc *BoardColumn) Save(db *gorm.DB) error {
	return db.Save(&bc).Error
}

type CreateBoardColumnArgs struct {
	BoardID string `json:"boardID" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

type UpdateBoardColumnArgs struct {
	ColumnID string  `json:"ColumnID" binding:"required"`
	Name     *string `json:"name"`
}

func (user *User) CreateBoardColumn(db *gorm.DB, args *CreateBoardColumnArgs) error {
	var bc BoardColumn

	bc.ID = uuid.New().String()
	bc.BoardID = args.BoardID
	bc.Name = args.Name

	err := bc.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "bc.Create",
			"subFunc": "bc.Create",
			"userID":  user.ID,
			"boardID": args.BoardID,
		}).Error(err)
		return err
	}

	return nil
}

func GetBoardColumns(db *gorm.DB, boardID string) (*[]BoardColumn, error) {
	var bc []BoardColumn

	err := db.Find(&bc, "board_id = ?", boardID).Error
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

func UpdateBoardColumn(db *gorm.DB, args *UpdateBoardColumnArgs) error {
	column, err := getBoardColumn(db, args.ColumnID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "UpdateBoardColumn",
			"subFunc":  "getBoardColumn",
			"columnID": args.ColumnID,
		}).Error(err)
		return err
	}

	if args.Name != nil {
		column.Name = *args.Name
	}

	err = column.Save(db)
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

func getBoardColumn(db *gorm.DB, columnID string) (*BoardColumn, error) {
	var column BoardColumn

	err := db.Find(&column, "id = ?", columnID).Error
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
