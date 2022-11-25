package createpdf

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
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
}

type JobWithZIP interface {
	UploadFile() error

	checkUploadedFile() bool

	readFile() error

	saveFile() error
}

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

func (f *FormFile) saveFile() error {
	var nameDir string = generateName()
	pathToTemporaryStorage := fmt.Sprintf("temporaryStorage/zips/%s", nameDir)

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
		NameDir: nameDir,
	}

	if err := createPDF.create(); err != nil {
		seterror.SetAppError("create.create()", err)
		return err
	}

	return nil
}

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

func (f *FormFile) checkUploadedFile() bool {
	if f.Header.Size > 2147483648 {
		return false
	}

	if f.Header.Header.Get("Content-Type") != "application/zip" {
		return false
	}

	return true
}

func (f *FormFile) UploadFile() error {
	var err error

	if err := f.Context.Request.ParseForm(); err != nil {
		seterror.SetAppError("context.Request.ParseForm", err)
		return err
	}

	f.File, f.Header, err = f.Context.Request.FormFile("zip")
	if err != nil {
		seterror.SetAppError("context.Request.FormFile", err)
		return err
	}

	if !f.checkUploadedFile() {
		return nil
	}

	if err := f.readFile(); err != nil {
		seterror.SetAppError("readFile", err)
		return err
	}

	return nil
}
