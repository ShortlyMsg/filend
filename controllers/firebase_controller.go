package controllers

import (
	"context"
	"encoding/json"
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

	bodyData := map[string]interface{}{
		"uploadedMB": progress.UploadedMB,
		"totalMB":    progress.TotalMB,
		"progress":   progress.Progress,
	}

	bodyJSON, err := json.Marshal(bodyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON oluşturma hatası: " + err.Error()})
		return
	}

	message := &messaging.Message{
		Topic: progress.Otp,
		Notification: &messaging.Notification{
			Title: progress.FileName,
			Body:  string(bodyJSON),
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase mesaj gönderme hatası: " + err.Error()})
		return
	}

	log.Printf("Mesaj başarıyla gönderildi. Response: %+v", response, progress)
	c.JSON(http.StatusOK, gin.H{"message": "Firebase bildirimi gönderildi", "response": response})
}
