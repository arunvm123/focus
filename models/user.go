package models

// User model
type User struct {
	ID          int     `json:"id" gorm:"primary_key"`
	Email       string  `json:"email" gorm:"unique;not null"`
	Name        string  `json:"name"`
	Password    string  `json:"password"`
	Verified    bool    `json:"verified"`
	ProfilePic  *string `json:"profilePic"`
	GoogleOauth bool    `json:"googleOauth"`
}

// UserProfile describes users profile
type UserProfile struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	Name       string  `json:"name"`
	ProfilePic *string `json:"profilePic"`
}

type SignUpArgs struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginArgs struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginWithGoogleArgs struct {
	Email   string  `json:"email" binding:"required,email"`
	Name    string  `json:"name" binding:"required"`
	Picture *string `json:"picture"`
}

type UpdateProfileArgs struct {
	// Email    *string `json:"email,omitempty"`
	Name       *string `json:"name,omitempty"`
	ProfilePic *string `json:"profilePic,omitempty"`
}

type UpdatePasswordArgs struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}
