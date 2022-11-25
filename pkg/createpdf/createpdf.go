package createpdf

import (
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html-to-pdf/pkg/seterror"
	"os"
)

type CreatePDF struct {
	NameDir string
}

type createPDF interface {
	create() (int64, error)
}

// getFileSize Получет размер pdf-файла
func getFileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		seterror.SetAppError("os.Open", err)
		return 0, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		seterror.SetAppError("file.Stat", err)
		return 0, err
	}

	return fileStat.Size(), nil
}

// create Создает pdf-файл
func (c *CreatePDF) create() (int64, error) {
	pdf, err := wkhtml.NewPDFGenerator()
	if err != nil {
		seterror.SetAppError("wkhtmltopdf.NewPDFGenerator", err)
		return 0, err
	}

	f, err := os.Open(fmt.Sprintf(`temporaryStorage/zips/%s/index.html`, c.NameDir))
	if err != nil {
		seterror.SetAppError("os.Open", err)
		return 0, err
	}

	page := wkhtml.NewPageReader(f)
	page.PageOptions.EnableLocalFileAccess.Set(true)
	pdf.AddPage(page)

	pdf.SetStderr(os.Stdout)

	if err := pdf.Create(); err != nil {
		seterror.SetAppError("pdf.Create()", err)
		return 0, err
	}

	if err := pdf.WriteFile(fmt.Sprintf("./temporaryStorage/pdfFiles/%s.pdf", c.NameDir)); err != nil {
		seterror.SetAppError("pdf.WriteFile", err)
		return 0, err
	}

	if err := os.RemoveAll(fmt.Sprintf("temporaryStorage/zips/%s", c.NameDir)); err != nil {
		seterror.SetAppError("os.RemoveAll", err)
		return 0, err
	}

	fileSize, err := getFileSize(fmt.Sprintf("./temporaryStorage/pdfFiles/%s.pdf", c.NameDir))
	if err != nil {
		seterror.SetAppError("getFileSize", err)
		return 0, err
	}

	return fileSize, nil
}
