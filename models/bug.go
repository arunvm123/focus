package models

type Bug struct {
	ID        int    `json:"id,gorm:"primary_key"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
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
