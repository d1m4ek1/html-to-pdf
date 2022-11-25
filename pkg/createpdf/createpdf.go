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
	create() error
}

func (c *CreatePDF) create() error {
	pdf, err := wkhtml.NewPDFGenerator()
	if err != nil {
		seterror.SetAppError("wkhtmltopdf.NewPDFGenerator", err)
		return err
	}

	f, err := os.Open(fmt.Sprintf(`temporaryStorage/zips/%s/index.html`, c.NameDir))
	if err != nil {
		seterror.SetAppError("os.Open", err)
		return err
	}

	page := wkhtml.NewPageReader(f)
	page.PageOptions.EnableLocalFileAccess.Set(true)
	pdf.AddPage(page)

	pdf.SetStderr(os.Stdout)

	if err := pdf.Create(); err != nil {
		seterror.SetAppError("pdf.Create()", err)
		return err
	}

	if err := pdf.WriteFile(fmt.Sprintf("./temporaryStorage/pdfFiles/%s.pdf", c.NameDir)); err != nil {
		seterror.SetAppError("pdf.WriteFile", err)
		return err
	}

	if err := os.RemoveAll(fmt.Sprintf("temporaryStorage/zips/%s", c.NameDir)); err != nil {
		seterror.SetAppError("os.RemoveAll", err)
		return err
	}

	return nil
}
