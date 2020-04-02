package emails

import (
	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

const (
	emailValidation string = "d-71b7ba0ed167416d8b07178b45cd2297"
	forgotPassword  string = "d-1ed04ad46478427d9970ba1b5a5a3033"
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
		"link": c.DomainURL + "forgotpassword/password/reset?token=" + token,
	}, forgotPassword)
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
