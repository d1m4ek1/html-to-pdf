package api

import (
	"github.com/gin-gonic/gin"
	"html-to-pdf/pkg/createpdf"
	"html-to-pdf/pkg/logging"
	"html-to-pdf/pkg/seterror"
	"net/http"
	"time"
)

/*
UploadFile Контроллер выполняющийся по запросу POST /api/upload. Загружает zip-файл и записывает содержимое в
временное хранилище zip-файлов
*/
func UploadFile() gin.HandlerFunc {
	return gin.HandlerFunc(func(context *gin.Context) {
		dateTime := time.Now()
		duration := time.Since(dateTime)

		var readFormFile createpdf.JobWithZIP = &createpdf.FormFile{
			Context: context,
		}
		fileName, fileSize, err := readFormFile.UploadFile()
		if err != nil {
			seterror.SetAppError("readFormFile.UploadFile()", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"successfully": false,
			})
			return
		}

		var logComplete logging.CreateLogging
		logComplete = logging.NewLogging(fileName, duration.String(), fileSize)

		if fileName != "" && fileSize != 0 {
			if err := logComplete.SetLogComplete(); err != nil {
				seterror.SetAppError("logComplete.SetLogComplete()", err)
				context.JSON(http.StatusInternalServerError, gin.H{
					"successfully": false,
				})
				return
			}
		} else {
			if err := logComplete.SetLogWarning("ZIP file not found"); err != nil {
				seterror.SetAppError("logComplete.SetLogComplete()", err)
				context.JSON(http.StatusInternalServerError, gin.H{
					"successfully": false,
				})
				return
			}
			context.JSON(http.StatusNotFound, gin.H{
				"successfully": false,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"successfully": true,
		})
	})
}
