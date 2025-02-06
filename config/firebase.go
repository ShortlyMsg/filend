package config

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App
var MessagingClient *messaging.Client

func InıtFirebase() {

	credJSON := os.Getenv("FIREBASE_CREDENTIALS")
	// JSON formatını []byte dizisine çevir
	opt := option.WithCredentialsJSON([]byte(credJSON))

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
