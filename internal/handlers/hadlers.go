package handlers

import (
	"github.com/gin-gonic/gin"
	"html-to-pdf/internal/controllers/api"
	"html-to-pdf/internal/controllers/templates"
)

// InitTemplateHandlers Инициализация запросов для загрузки HTML-шаблонов
func InitTemplateHandlers(rtr *gin.Engine) {
	index := rtr.Group("/")
	{
		index.GET("/", templates.HomeTemplate())
	}
}

// InitAPIHandlers Инициализация api запросов
func InitAPIHandlers(rtr *gin.Engine) {
	upload := rtr.Group("/api")
	{
		upload.POST("/upload", api.UploadFile())
	}
}
