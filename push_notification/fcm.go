package push

import (
	"context"
	"log"

	"firebase.google.com/go/messaging"
)

type Payload struct {
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data"`
	ClickAction string            `json:"clickAction"`
}

func SendPushNotification(client *messaging.Client, registrationTokens []string, payload *Payload) error {
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

	_, err := client.SendMulticast(context.Background(), notification)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
