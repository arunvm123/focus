package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// List model
type List struct {
	ID        int    `json:"id" gorm:"primary_key"`
	UserID    int    `json:"userId"`
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

// CreateListArgs defines the args for create list api
type CreateListArgs struct {
	Heading string `json:"heading"`
}

func getListOfUser(db *gorm.DB, userID int) (*List, error) {
	var list List

	err := db.Find(&list, "user_id = ? AND archived = false", userID).Error
	if err != nil {
		log.Printf("Error when fetching list\n%v", err)
		return nil, err
	}

	return &list, nil
}

func (user *User) CreateList(db *gorm.DB, args *CreateListArgs) error {
	list := List{
		UserID:    user.ID,
		Archived:  false,
		CreatedAt: time.Now().Unix(),
		Heading:   args.Heading,
	}

	err := list.Create(db)
	if err != nil {
		log.Printf("Error when creating list\n%v", err)
		return err
	}

	return nil
}
