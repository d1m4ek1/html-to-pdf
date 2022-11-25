package handlers

import (
	"github.com/gin-gonic/gin"
	"html-to-pdf/internal/controllers/api"
	"html-to-pdf/internal/controllers/templates"
)

func InitTemplateHandlers(rtr *gin.Engine) {
	index := rtr.Group("/")
	{
		index.GET("/", templates.HomeTemplate())
	}
}

func InitAPIHandlers(rtr *gin.Engine) {
	upload := rtr.Group("/api")
	{
		upload.POST("/upload", api.UploadFile())
	}
}
