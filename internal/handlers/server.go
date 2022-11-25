package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitStaticFiles(rtr *gin.Engine) {
	rtr.StaticFS("/static/", http.Dir("./ui/static/"))
}

func InitHTMLTemplates(rtr *gin.Engine) {
	rtr.LoadHTMLFiles("./ui/index.html")
}

func InitGin() {
	gin.SetMode(gin.ReleaseMode)
	rtr := gin.Default()

	InitStaticFiles(rtr)

	InitHTMLTemplates(rtr)

	InitTemplateHandlers(rtr)

	InitAPIHandlers(rtr)

	rtr.Run(":3000")
}
