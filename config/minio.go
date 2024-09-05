package config

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	endpoint := "localhost:9000"
	accessKeyID := "Zjr7F3ddORJnFMq3avF5"
	secretAccessKey := "X1uEr7avB00CoaiQ8sF6iTpCTEQeXtADJNrFPInn"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	MinioClient = client
	log.Println("Connected to MinIO")
}
