package mysql

import (
	"time"

	"github.com/arunvm/focus/models"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateBug(args *models.CreateBugArgs, user *models.User) error {
	bug := models.Bug{
		Status:    models.TODO,
		CreatedAt: time.Now().Unix(),
		Info:      args.Info,
		Title:     args.Title,
		UserID:    user.ID,
	}

	err := db.Client.Create(&bug).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "CreateBug",
			"subFunc": "bug.Create",
			"userID":  user.ID,
			"args":    *args,
		}).Error(err)
		return err
	}

	return nil
}

func (db *Mysql) GetBugs() (*[]models.BugInfo, error) {
	var bugs []models.BugInfo

	err := db.Client.Table("bugs").Joins("JOIN users on bugs.user_id = users.id").
		Select("bugs.*,users.profile_pic,users.name").
		Find(&bugs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "GetBugs",
			"subFunc": "retrieving bug info",
		}).Error(err)
		return nil, err
	}

	return &bugs, nil
}

func (db *Mysql) UpdateBug(args *models.UpdateBugArgs) error {
	var bug models.Bug

	err := db.Client.Find(&bug, "id = ?", args.ID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "UpdateBug",
			"info":  "retrieving bug with id",
			"bugID": args.ID,
		}).Error(err)
		return err
	}

	bug.Status = args.Status
	err = db.Client.Save(bug).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateBug",
			"subFunc": "bug.Save",
			"bugID":   args.ID,
		}).Error(err)
		return err
	}

	return nil
}
