package fcm

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	push "github.com/arunvm/travail-backend/push_notification"
	"google.golang.org/api/option"
)

type FCM struct {
	Client *messaging.Client
}

func New(serviceAccountKeyPath string) (*FCM, error) {
	firebaseApp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return nil, err
	}

	var fcm FCM

	fcm.Client, err = firebaseApp.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &fcm, nil
}

func (fcm *FCM) SendPushNotification(registrationTokens []string, payload *push.Payload) error {
	notification := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: "https://i.imgur.com/mGOXXII.png",
		},
		Data: payload.Data,
		Webpush: &messaging.WebpushConfig{
			Data: payload.Data,
			Notification: &messaging.WebpushNotification{
				Title:   payload.Title,
				Body:    payload.Body,
				Icon:    "https://i.imgur.com/mGOXXII.png",
				Vibrate: []int{200, 100, 200},
			},
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Color:                 "#4C51BF",
				ClickAction:           payload.ClickAction,
				DefaultSound:          true,
				DefaultVibrateTimings: true,
				DefaultLightSettings:  true,
			},
			Data: payload.Data,
		},
		Tokens: registrationTokens,
	}

	_, err := fcm.Client.SendMulticast(context.Background(), notification)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
