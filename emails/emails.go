package emails

import (
	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/models"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

const (
	emailValidation    string = "d-71b7ba0ed167416d8b07178b45cd2297"
	forgotPassword     string = "d-1ed04ad46478427d9970ba1b5a5a3033"
	organisationInvite string = "d-fafbec3bf89944f5ba7f991ca19d81ab"
)

func SendValidationEmail(emailClient *sendgrid.Client, user *models.User, token string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Name:    user.Name,
		Address: user.Email,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendValidationEmail",
			"subFunc": "config.GetConfig",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return sendEmail(emailClient, to, map[string]interface{}{
		"name":            user.Name,
		"validation_link": c.DomainURL + "verify/module?token=" + token,
	}, emailValidation)
}

func SendForgotPasswordEmail(emailClient *sendgrid.Client, user *models.User, token string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Name:    user.Name,
		Address: user.Email,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendForgotPasswordnEmail",
			"subFunc": "config.GetConfig",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return sendEmail(emailClient, to, map[string]interface{}{
		"name": user.Name,
		"link": c.DomainURL + "forgot/module?token=" + token,
	}, forgotPassword)
}

func SendOrganisationInvite(emailClient *sendgrid.Client, db *gorm.DB, adminName string, invite *models.OrganisationInvitation) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Address: invite.Email,
	})

	orgName, err := models.GetOrganisationName(db, invite.OrganisationID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":          "SendInviteToOrganisation",
			"subFunc":       "models.GetOrganisationName",
			"organistionID": invite.OrganisationID,
			"email":         invite.Email,
		}).Error(err)
		return err
	}

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":          "SendInviteToOrganisation",
			"subFunc":       "config.GetConfig",
			"organistionID": invite.OrganisationID,
			"email":         invite.Email,
		}).Error(err)
		return err
	}

	return sendEmail(emailClient, to, map[string]interface{}{
		"admin_name":        adminName,
		"organisation_name": orgName,
		"invite_link":       c.DomainURL + "organisation/invite/accept?token=" + invite.Token,
	}, organisationInvite)
}

func sendEmail(emailCLient *sendgrid.Client, to []*mail.Email, templateData map[string]interface{}, templateID string) error {
	personalizations := []*mail.Personalization{}
	personalizations = append(personalizations, &mail.Personalization{
		To:                  to,
		DynamicTemplateData: templateData,
	})

	email := mail.NewV3Mail()

	email.SetFrom(&mail.Email{
		Address: "info@travail.com",
		Name:    "Travail",
	})
	email.SetTemplateID(templateID)
	email.AddPersonalizations(personalizations...)

	resp, err := emailCLient.Send(email)
	if err != nil {
		log.WithFields(log.Fields{
			"func":       "sendEmail",
			"subFunc":    "emailCLient.Send",
			"templateID": templateID,
		}).Error(err)
		return err
	}

	log.Println(resp)

	return nil
}
