package config

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App
var MessagingClient *messaging.Client

func InitFirebase() {

	opt := option.WithCredentialsFile("./firebase-config-be.json")

	// Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase başlatma hatası: %v", err)
	}

	// Firebase Messaging Client
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Firebase Messaging başlatma hatası: %v", err)
	}

	FireBaseApp = app
	MessagingClient = client

	fmt.Println("🔥 Firebase başarıyla başlatıldı!")
}
