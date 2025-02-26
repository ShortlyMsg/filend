package routes

import (
	"filend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/upload", controllers.UploadFile)
	router.POST("/download", controllers.DownloadFile)
	router.POST("/getAllFiles", controllers.GetAllFiles)
	router.POST("/checkFileHash", controllers.CheckFileHash)
	router.POST("/generateOtp", controllers.GenerateOTP)
	router.POST("/sendProgress", controllers.SendUploadProgress)

	return router
}
