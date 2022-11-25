package handlers

import (
	"github.com/gin-gonic/gin"
	"html-to-pdf/pkg/seterror"
	"log"
	"net/http"
)

// InitStaticFiles Инициализация статичных файлов и загрузка на сайт по запросу
func InitStaticFiles(rtr *gin.Engine) {
	rtr.StaticFS("/static/", http.Dir("./ui/static/"))
}

// InitHTMLTemplates Инициализация HTML шаблонов
func InitHTMLTemplates(rtr *gin.Engine) {
	rtr.LoadHTMLFiles("./ui/index.html")
}

// InitGin Инициализация фреймворка gin-gonic и создание сервера
func InitGin() error {
	gin.SetMode(gin.ReleaseMode)

	log.Println("Init gin-gonic")
	rtr := gin.Default()

	log.Println("Init static files")
	InitStaticFiles(rtr)

	log.Println("Init HTML templates")
	InitHTMLTemplates(rtr)

	log.Println("Init template handlers")
	InitTemplateHandlers(rtr)

	log.Println("Init API handlers")
	InitAPIHandlers(rtr)

	log.Println("Run server, url - localhost:3000")
	if err := rtr.Run(":3000"); err != nil {
		seterror.SetAppError("rtr.Run", err)
		return err
	}

	return nil
}
