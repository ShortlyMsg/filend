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
	router.GET("/download/:otp", controllers.DownloadFile)
	router.POST("/getAllFiles", controllers.GetAllFiles)

	return router
}
