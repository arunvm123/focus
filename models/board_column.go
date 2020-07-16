package models

type BoardColumn struct {
	ID      string `json:"id" gorm:"primary_key;auto_increment:false"`
	BoardID string `json:"boardID"`
	Name    string `json:"name"`
}

type CreateBoardColumnArgs struct {
	BoardID string `json:"boardID"`
	Name    string `json:"name" binding:"required"`
}

type UpdateBoardColumnArgs struct {
	BoardID  string  `json:"boardID"`
	ColumnID string  `json:"ColumnID" binding:"required"`
	Name     *string `json:"name"`
}
