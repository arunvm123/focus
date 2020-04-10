package models

import (
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Bug struct {
	ID        int    `json:"id,gorm:"primary_key"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

func (t *Bug) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

func (t *Bug) Save(db *gorm.DB) error {
	return db.Save(&t).Error
}

const (
	TODO   = 1
	FIXING = 2
	FIXED  = 3
)

type CreateBugArgs struct {
	Title string `json:"title" binding:"required"`
	Info  string `json:"info" binding:"required"`
}

type BugInfo struct {
	Bug
	Name       string  `json:"name"`
	ProfilePic *string `json:"profile_pic"`
}

type UpdateBugArgs struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required,eq=2|eq=3|eq=1"`
}

func (user *User) CreateBug(db *gorm.DB, args *CreateBugArgs) error {
	bug := Bug{
		Status:    TODO,
		CreatedAt: time.Now().Unix(),
		Info:      args.Info,
		Title:     args.Title,
		UserID:    user.ID,
	}

	err := bug.Create(db)
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

func (admin *User) GetBugs(db *gorm.DB) (*[]BugInfo, error) {
	var bugs []BugInfo

	err := db.Table("bugs").Joins("JOIN users on bugs.user_id = users.id").
		Select("bugs.*,users.profile_pic,users.name").
		Find(&bugs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "GetBugs",
			"subFunc": "retrieving bug info",
			"adminID": admin.ID,
		}).Error(err)
		return nil, err
	}

	return &bugs, nil
}

func (admin *User) UpdateBug(db *gorm.DB, args *UpdateBugArgs) error {
	var bug Bug

	err := db.Find(&bug, "id = ?", args.ID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateBug",
			"info":    "retrieving bug with id",
			"adminID": admin.ID,
			"bugID":   args.ID,
		}).Error(err)
		return err
	}

	bug.Status = args.Status
	err = bug.Save(db)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateBug",
			"subFunc": "bug.Save",
			"adminID": admin.ID,
			"bugID":   args.ID,
		}).Error(err)
		return err
	}

	return nil
}
