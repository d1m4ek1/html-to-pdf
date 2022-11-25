package api

import (
	"github.com/gin-gonic/gin"
	"html-to-pdf/pkg/createpdf"
	"html-to-pdf/pkg/seterror"
	"net/http"
)

func UploadFile() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		var readFormFile createpdf.JobWithZIP = &createpdf.FormFile{
			Context: context,
		}
		if err := readFormFile.UploadFile(); err != nil {
			seterror.SetAppError("readFormFile.UploadFile()", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"successfully": false,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"successfully": true,
		})
	})
}
