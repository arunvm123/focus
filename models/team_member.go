package models

type TeamMember struct {
	TeamID   string `json:"teamID" gorm:"primary_key;auto_increment:false"`
	UserID   int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	JoinedAt int64  `json:"joinedAt"`
}

type AddTeamMemberArgs struct {
	TeamID string `json:"-"`
	UserID int    `json:"userID" binding:"required"`
}

type TeamMemberInfo struct {
	TeamID     string  `json:"teamID"`
	UserID     int     `json:"userID"`
	Name       string  `json:"name"`
	ProfilePic *string `json:"profile_pic"`
}
