package controllers

import (
	"context"
	"filend/config"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
)

type UploadProgress struct {
	Otp      string `json:"otp"`
	FileName string `json:"fileName"`
	Progress int    `json:"progress"`
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
		Data: map[string]string{
			"fileName": progress.FileName,
			"progress": strconv.Itoa(progress.Progress),
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

type TopicSubscriptionRequest struct {
	Token string `json:"token"`
	Topic string `json:"topic"`
}

func SubscribeTokenToTopic(c *gin.Context) {
	var req TopicSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz JSON formatı"})
		return
	}

	client := config.MessagingClient

	// Token'ı belirli topic'e abone et
	response, err := client.SubscribeToTopic(context.Background(), []string{req.Token}, req.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Abonelik işlemi başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Cihaz topic'e abone edildi",
		"response": response,
	})
}
