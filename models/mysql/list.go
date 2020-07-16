package mysql

import (
	"errors"
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateList(args *models.CreateListArgs, user *models.User) (*models.List, error) {
	team, err := db.getTeamFromID(args.TeamID)
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
		return nil, models.UserNotAdminOfTeam
	}

	list := models.List{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Heading:   args.Heading,
		TeamID:    args.TeamID,
	}

	err = db.Client.Create(list).Error
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
func (db *Mysql) GetLists(args *models.GetListsArgs, user *models.User) (*[]models.ListInfo, error) {
	var lists []models.ListInfo

	team, err := db.getTeamFromID(args.TeamID)
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

	err = db.Client.Table("lists").Joins("LEFT JOIN tasks on lists.id = tasks.list_id").
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
func (db *Mysql) UpdateList(args *models.UpdateListArgs, user *models.User) error {
	list, err := db.getList(args.ID)
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

	err = db.Client.Save(list).Error
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

func (db *Mysql) getList(listID string) (*models.List, error) {
	var list models.List

	err := db.Client.Find(&list, "archived = false AND id = ?", listID).Error
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
