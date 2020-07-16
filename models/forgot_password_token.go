package models

type ForgotPasswordToken struct {
	UserID    int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
