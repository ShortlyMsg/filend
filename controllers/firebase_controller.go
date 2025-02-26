package controllers

import (
	"context"
	"filend/config"
	"log"
	"net/http"

	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
)

type UploadProgress struct {
	Otp        string `json:"otp"`
	FileName   string `json:"fileName"`
	UploadedMB string `json:"uploadedMB"`
	TotalMB    string `json:"totalMB"`
	Progress   int    `json:"progress"`
}

func SendUploadProgress(c *gin.Context) {
	var progress UploadProgress

	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz JSON formatı"})
		return
	}

	client := config.MessagingClient

	message := &messaging.Message{
		Topic: progress.Otp,
		Notification: &messaging.Notification{
			Title: progress.FileName,
			Body:  progress.UploadedMB + "MB / " + progress.TotalMB + "MB (" + string(progress.Progress) + "%)",
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase mesaj gönderme hatası: " + err.Error()})
		return
	}

	log.Printf("Mesaj başarıyla gönderildi. Response: %+v", response)
	c.JSON(http.StatusOK, gin.H{"message": "Firebase bildirimi gönderildi", "response": response})
}
