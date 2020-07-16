package models

type EmailValidateToken struct {
	UserID     int    `json:"userID" gorm:"primary_key;auto_increment:false"`
	Token      string `json:"token" gorm:"primary_key;auto_increment:false"`
	CreatedAt  int64  `json:"createdAt"`
	ExpiresAt  int64  `json:"expiresAt"`
	Invalidate bool   `json:"invalidate"` // Token invalidates when new token is generated, when a new email verification is rewuested
}

type ValidateEmailArgs struct {
	Token string `json:"token" binding:"required"`
}
