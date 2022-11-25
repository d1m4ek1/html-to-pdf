package createpdf

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"html-to-pdf/pkg/logging"
	"html-to-pdf/pkg/seterror"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"strings"
)

type FormFile struct {
	File        multipart.File
	ReadiedFile *zip.Reader
	Header      *multipart.FileHeader
	Context     *gin.Context

	NewName         string
	CreatedFileSize int64
}

type JobWithZIP interface {
	UploadFile() (string, int64, error)

	checkUploadedFile(createLogging logging.CreateLogging) bool

	readFile() error

	saveFile() error
}

// generateName Генерирует имя pdf-файла
func generateName() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var characters = 12

	b := make([]rune, characters)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			seterror.SetAppError("rand.Int", err)
			return ""
		}

		b[i] = letterRunes[n.Int64()]
	}
	return string(b)
}

// validExtension Проверяет допустимые расширения файлов внутри zip-файла
func validExtension(filename, filetype string) bool {
	switch filetype {
	case "files":
		if strings.Contains(filename, ".html") || strings.Contains(filename, ".css") {
			return true
		}
		return false

	case "pictures":
		if strings.Contains(filename, ".jpg") || strings.Contains(filename, ".jpeg") ||
			strings.Contains(filename, ".png") || strings.Contains(filename, ".gif") {
			return true
		}
		return false
	}

	return false
}

// saveFile Сохраняет файлы zip-файла в временное хранилище файлов
func (f *FormFile) saveFile() error {
	f.NewName = generateName()
	pathToTemporaryStorage := fmt.Sprintf("temporaryStorage/zips/%s", f.NewName)

	if err := os.MkdirAll(pathToTemporaryStorage, 0777); err != nil {
		seterror.SetAppError("os.MkdirAll", err)
		return err
	}

	for _, file := range f.ReadiedFile.File {
		if validExtension(file.Name, "files") {
			loadFiles(file, pathToTemporaryStorage)
		}
		if validExtension(file.Name, "pictures") {
			loadPictures(file, pathToTemporaryStorage)
		}
	}

	var createPDF createPDF = &CreatePDF{
		NameDir: f.NewName,
	}

	fileSize, err := createPDF.create()
	if err != nil {
		seterror.SetAppError("create.create()", err)
		return err
	}

	f.CreatedFileSize = fileSize

	return nil
}

// readFile Читает и получает файлы, zip-файла
func (f *FormFile) readFile() error {
	var err error
	buffer := new(bytes.Buffer)

	fileSize, err := io.Copy(buffer, f.File)
	if err != nil {
		seterror.SetAppError("io.Copy", err)
		return err
	}

	f.ReadiedFile, err = zip.NewReader(bytes.NewReader(buffer.Bytes()), fileSize)
	if err != nil {
		seterror.SetAppError("zip.NewReader", err)
		return err
	}

	if err := f.saveFile(); err != nil {
		seterror.SetAppError("f.openFile", err)
		return err
	}

	return nil
}

// checkUploadedFile Проверяет загруженный файл на достоверность типа и размера
func (f *FormFile) checkUploadedFile(createLogging logging.CreateLogging) bool {
	if f.Header.Size > 2147483648 {
		createLogging.SetLogWarning("The file size exceeds the allowed 2 gigabytes")
		return false
	}

	if f.Header.Header.Get("Content-Type") != "application/zip" {
		createLogging.SetLogWarning("The file type is not a zip file")
		return false
	}

	return true
}

// UploadFile Получает файлы из запроса, затем читает и сохраняет файлы
func (f *FormFile) UploadFile() (string, int64, error) {
	var err error

	var createLogging logging.CreateLogging = logging.NewLogging("", "", 0)

	if err := f.Context.Request.ParseForm(); err != nil {
		seterror.SetAppError("context.Request.ParseForm", err)
		return "", 0, err
	}

	f.File, f.Header, err = f.Context.Request.FormFile("zip")
	if f.File == nil {
		return "", 0, nil
	}
	if err != nil {
		seterror.SetAppError("context.Request.FormFile", err)
		return "", 0, err
	}

	if !f.checkUploadedFile(createLogging) {
		return "", 0, nil
	}

	if err := f.readFile(); err != nil {
		seterror.SetAppError("readFile", err)
		return "", 0, err
	}

	return f.NewName + ".pdf", f.CreatedFileSize, nil
}
