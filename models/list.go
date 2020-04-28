package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// List model
type List struct {
	ID        string `json:"id" gorm:"primary_key"`
	UserID    int    `json:"userId"`
	TeamID    string `json:"teamID"`
	Heading   string `json:"heading"`
	CreatedAt int64  `json:"createdAt"`
	Archived  bool   `json:"archived"`
}

// Create is a helper function to create a new list
func (l *List) Create(db *gorm.DB) error {
	return db.Create(&l).Error
}

// Save is a helper function to update existing list
func (l *List) Save(db *gorm.DB) error {
	return db.Save(&l).Error
}

type ListInfo struct {
	List
	PendingTasks   int  `json:"pendingTasks"`
	CompletedTasks int  `json:"completedTasks"`
	Active         bool `json:"active"` //Additional info to help with UI
}

// CreateListArgs defines the args for create list api
type CreateListArgs struct {
	TeamID  string `json:"teamID" binding:"required"`
	Heading string `json:"heading"`
}

type GetListsArgs struct {
	TeamID string `json:"teamID" binding:"required"`
}

// UpdateListArgs defines the args for update list api
type UpdateListArgs struct {
	ID       string  `json:"id" binding:"required"`
	Heading  *string `json:"heading,omitempty"`
	Archived *bool   `json:"archived,omitempty"`
}

// func getListOfUser(db *gorm.DB, userID int) (*List, error) {
// 	var list List

// 	err := db.Find(&list, "user_id = ? AND archived = false", userID).Error
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"func":   "getListOfUser",
// 			"info":   "retrieving all lists of user",
// 			"userID": userID,
// 		}).Error(err)
// 		return nil, err
// 	}

// 	return &list, nil
// }

func (user *User) CreateList(db *gorm.DB, args *CreateListArgs) (*List, error) {
	team, err := getTeamFromID(db, args.TeamID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateList",
			"subFunc": "getTeamFromID",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	if team.AdminID != user.ID {
		return nil, errors.New("User not admin of team")
	}

	list := List{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Heading:   args.Heading,
		TeamID:    args.TeamID,
	}

	err = list.Create(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateList",
			"subFunc": "list.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	return &list, nil
}

// GetLists returns all lists of the user
func (user *User) GetLists(db *gorm.DB, args *GetListsArgs) (*[]ListInfo, error) {
	var lists []ListInfo

	team, err := getTeamFromID(db, args.TeamID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "GetLists",
			"subFunc": "getTeamFromID",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return nil, err
	}

	if team.AdminID != user.ID {
		return nil, errors.New("User not admin of team")
	}

	err = db.Table("lists").Joins("LEFT JOIN tasks on lists.id = tasks.list_id").
		Select("lists.*,"+
			"sum(case when complete = true AND tasks.archived = false then 1 else 0 end) as completed_tasks,"+
			"sum(case when complete = false AND tasks.archived = false then 1 else 0 end) as pending_tasks").
		Where("lists.archived = false AND lists.user_id = ? AND lists.team_id = ?", user.ID, team.ID).
		Group("lists.id").
		Find(&lists).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetLists",
			"info":   "retrieving list info",
			"userID": user.ID,
		}).Error(err)
		return nil, err
	}

	return &lists, nil
}

// UpdateList updates list info
func (user *User) UpdateList(db *gorm.DB, args *UpdateListArgs) error {
	list, err := getList(db, args.ID)
	if err != nil {
		log.Printf("Error while getting list\n%v", err)
		return err
	}

	if list.UserID != user.ID {
		return errors.New("Not user's list")
	}

	if args.Heading != nil {
		list.Heading = *args.Heading
	}
	if args.Archived != nil {
		list.Archived = *args.Archived
	}

	err = list.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateList",
			"subFunc": "list.Save",
			"userID":  user.ID,
			"listID":  args.ID,
		}).Error(err)
		return err
	}

	return nil
}

func getList(db *gorm.DB, listID string) (*List, error) {
	var list List

	err := db.Find(&list, "archived = false AND id = ?", listID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "getList",
			"info":   "retrieving list info",
			"listID": listID,
		}).Error(err)
		return nil, err
	}

	return &list, nil
}
