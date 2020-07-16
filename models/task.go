package models

// Task model
type Task struct {
	ID        string `json:"id" gorm:"primary_key"`
	ListID    string `json:"listID"`
	Info      string `json:"info"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
	ExpiresAt *int64 `json:"expiresAt"`
	Complete  bool   `json:"complete"`
	Order     int    `json:"order"`
	Archived  bool   `json:"archived"`
}

type CreateTaskArgs struct {
	ListID    string `json:"listID" binding:"required"`
	Info      string `json:"info"`
	Order     int    `json:"order"`
	ExpiresAt *int64 `json:"expiresAt"`
}

type GetTasksArgs struct {
	ListID   string `json:"listID" binding:"required"`
	Complete *bool  `json:"complete"`
}

type UpdateTaskArgs struct {
	ID        string  `json:"id" binding:"required"`
	Info      *string `json:"info"`
	Order     *int    `json:"order"`
	Complete  *bool   `json:"complete"`
	Archived  *bool   `json:"archived"`
	ExpiresAt *int64  `json:"expiresAt"`
}

type TaskInfo struct {
	Task
	Heading string `json:"heading"`
	UserID  int    `json:"userID"`
}

type DeleteTasksArgs struct {
	TaskIDs []string `json:"taskIDs" binding:"required"`
}
