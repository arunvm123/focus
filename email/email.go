package email

import (
	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
)

type Email interface {
	SendValidationEmail(user *models.User, token string) error
	SendForgotPasswordEmail(user *models.User, token string) error
	SendOrganisationInvite(db *gorm.DB, adminName string, invite *models.OrganisationInvitation) error
}
