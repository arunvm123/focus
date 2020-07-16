package models

type FCMNotificationToken struct {
	UserID int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	Token  string `json:"token" gorm:"primary_key"`
}

type AddNotificationTokenArgs struct {
	Token string `json:"token" binding:"required"`
}
