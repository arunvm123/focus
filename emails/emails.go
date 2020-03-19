package emails

import (
	"log"

	"github.com/arunvm/travail-backend/config"
	"github.com/arunvm/travail-backend/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	emailValidation string = "d-71b7ba0ed167416d8b07178b45cd2297"
)

func SendValidationEmail(emailCLient *sendgrid.Client, user *models.User, token string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Name:    user.Name,
		Address: user.Email,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.Printf("Error reading config\n%v", err)
		return err
	}

	return sendEmail(emailCLient, to, map[string]interface{}{
		"validation_link": c.DomainURL + "verify?token=" + token,
	}, emailValidation)
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
		log.Println(err)
		return err
	}

	log.Println(resp)

	return nil
}
