package templates

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeTemplate() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		context.HTML(http.StatusOK, "index", nil)
	})
}
