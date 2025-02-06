package services

import (
	"context"
	"filend/config"
	"fmt"

	"firebase.google.com/go/messaging"
)

func SendPushNotification(token string, title string, body string) error {
	client := config.MessagingClient

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("🔥 Mesaj gönderme hatası: %v", err)
	}

	fmt.Println("Mesaj başarıyla gönderildi:", response)
	return nil
}
