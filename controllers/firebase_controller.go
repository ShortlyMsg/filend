package controllers

import (
	"fmt"
	"net/http"
)

func NotifyUploadStart(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		FileName string `json:"fileName"`
		FileSize string `json:"fileSize"`
		FileHash string `json:"fileHash"`
		Token    string `json:"token"`
	}

	fmt.Fprintf(w, "Upload başlatıldı: %s", requestData.FileName)
}

func NotifyDownloadStart(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Download başlatıldı: %s")
}
