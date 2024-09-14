package config

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	endpoint := "localhost:9000"
	accessKeyID := os.Getenv("MINIO_ACC_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACC_KEY")
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
