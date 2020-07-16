package models

type Board struct {
	ID        string `json:"id" gorm:"primary_key;auto_increment:false"`
	AdminID   int    `json:"adminID"`
	TeamID    string `json:"teamID"`
	Title     string `json:"title"`
	CreatedOn int64  `json:"createdOn"`
	Archived  bool   `json:"archived"`
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

type GetBoardsArgs struct {
	TeamID string `json:"teamID"`
}
