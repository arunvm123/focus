package models

type ColumnCard struct {
	ID          string `json:"id" gorm:"primary_key"`
	ColumnID    string `json:"columnID"`
	Heading     string `json:"heading"`
	Description string `json:"description" gorm:"size:3000"`
	AssignedTo  *int   `json:"assignedTo"`
	AssignedBy  *int   `json:"AssignedBy"`
	AssignedOn  *int64 `json:"assignedOn"`
}

type CreateColumnCardArgs struct {
	ColumnID    string `json:"columnID"`
	Heading     string `json:"heading" binding:"required"`
	Description string `json:"description"`
	AssignedTo  *int   `json:"assignedTo"`
	AssignedBy  *int   `json:"AssignedBy"`
}

type UpdateColumnCardArgs struct {
	CardID      string  `json:"cardID" binding:"required"`
	ColumnID    string  `json:"columnID"`
	Heading     *string `json:"heading"`
	Description *string `json:"description" gorm:"size:3000"`
	AssignedTo  *int    `json:"assignedTo"`
}
