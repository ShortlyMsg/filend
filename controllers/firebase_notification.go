package controllers

import (
	"context"
	"filend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
