package templates

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomeTemplate Загружает домашнюю страницу
func HomeTemplate() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		context.HTML(http.StatusOK, "index", nil)
	})
}
