package createpdf

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"html-to-pdf/pkg/seterror"
	"io"
	"mime/multipart"
	"strings"
)

type FormFile struct {
	File        multipart.File
	ReadiedFile *zip.Reader
	Header      *multipart.FileHeader
	Context     *gin.Context
}

type IFormFile interface {
	UploadFile() error

	checkUploadedFile() bool

	readFile() error

	openFile() error
}

func (f *FormFile) openFile() error {
	for _, file := range f.ReadiedFile.File {
		if strings.Contains(file.Name, ".html") {
			contentReadCloser, err := file.Open()
			if err != nil {
				seterror.SetAppError("file.Open", err)
				return err
			}

			buffer := new(strings.Builder)

			if _, err := io.Copy(buffer, contentReadCloser); err != nil {
				seterror.SetAppError("io.Copy(buffer, contentReadCloser)", err)
				return err
			}

			fmt.Println(buffer.String())
		}
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

	if err := f.openFile(); err != nil {
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
