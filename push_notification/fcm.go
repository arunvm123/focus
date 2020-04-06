package push

import (
	"context"
	"log"

	"firebase.google.com/go/messaging"
)

func sendPushNotification(client messaging.Client, registrationTokens []string, message string) error {
	notification := &messaging.MulticastMessage{
		Data: map[string]string{
			"info": message,
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
