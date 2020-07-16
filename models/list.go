package models

import (
	"errors"
)

var (
	UserNotAdminOfTeam = errors.New("User not admin of team")
)

// List model
type List struct {
	ID        string `json:"id" gorm:"primary_key"`
	UserID    int    `json:"userID"`
	TeamID    string `json:"teamID"`
	Heading   string `json:"heading"`
	CreatedAt int64  `json:"createdAt"`
	Archived  bool   `json:"archived"`
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
