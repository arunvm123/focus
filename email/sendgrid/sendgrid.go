package sendgrid

import (
	"github.com/arunvm/travail-backend/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

const (
	emailValidation    string = "d-71b7ba0ed167416d8b07178b45cd2297"
	forgotPassword     string = "d-1ed04ad46478427d9970ba1b5a5a3033"
	organisationInvite string = "d-fafbec3bf89944f5ba7f991ca19d81ab"
)

type Sendgrid struct {
	Client *sendgrid.Client
}

func New(sendgridKey string) *Sendgrid {
	client := sendgrid.NewSendClient(sendgridKey)

	return &Sendgrid{
		Client: client,
	}
}

func (sendgrid *Sendgrid) SendValidationEmail(name, email string, token string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Name:    name,
		Address: email,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendValidationEmail",
			"subFunc": "config.GetConfig",
		}).Error(err)
		return err
	}

	return sendEmail(sendgrid.Client, to, map[string]interface{}{
		"name":            name,
		"validation_link": c.DomainURL + "verify/module?token=" + token,
	}, emailValidation)
}

func (sendgrid *Sendgrid) SendForgotPasswordEmail(name, email string, token string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Name:    name,
		Address: email,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendForgotPasswordnEmail",
			"subFunc": "config.GetConfig",
		}).Error(err)
		return err
	}

	return sendEmail(sendgrid.Client, to, map[string]interface{}{
		"name": name,
		"link": c.DomainURL + "forgot/module?token=" + token,
	}, forgotPassword)
}

func (sendgrid *Sendgrid) SendOrganisationInvite(adminName, inviteEmail, token, orgName string) error {
	to := []*mail.Email{}
	to = append(to, &mail.Email{
		Address: inviteEmail,
	})

	c, err := config.GetConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendInviteToOrganisation",
			"subFunc": "config.GetConfig",
		}).Error(err)
		return err
	}

	return sendEmail(sendgrid.Client, to, map[string]interface{}{
		"admin_name":        adminName,
		"organisation_name": orgName,
		"invite_link":       c.DomainURL + "organisation/invite/accept?token=" + token,
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
