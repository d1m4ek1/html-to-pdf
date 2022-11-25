package createpdf

import (
	"archive/zip"
	"fmt"
	"html-to-pdf/pkg/seterror"
	"io"
	"os"
	"regexp"
	"strings"
)

func containAndRegexpHTML(filename, content, pathToTemporaryStorage string) string {
	if strings.Contains(filename, ".html") {
		path, _ := os.Getwd()

		regexpHREF := regexp.MustCompile(`href="|href="/|href="\./`)
		regexpSRC := regexp.MustCompile(`src="|src="/|src=\./`)

		content = regexpHREF.ReplaceAllString(content, fmt.Sprintf(`href="%s/%s/`,
			path, pathToTemporaryStorage))

		content = regexpSRC.ReplaceAllString(content, fmt.Sprintf(`src="%s/%s/`,
			path, pathToTemporaryStorage))
	}

	return content
}

func loadFiles(file *zip.File, pathToTemporaryStorage string) error {
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

	content := containAndRegexpHTML(file.Name, buffer.String(), pathToTemporaryStorage)

	fileCreated, err := os.Create(pathToTemporaryStorage + "/" + file.FileInfo().Name())
	if err != nil {
		seterror.SetAppError("os.Create", err)
		return err
	}

	if err := fileCreated.Close(); err != nil {
		seterror.SetAppError("fileCreated.Close()", err)
	}

	if err := os.WriteFile(pathToTemporaryStorage+"/"+file.FileInfo().Name(), []byte(content), 0666); err != nil {
		seterror.SetAppError("os.WriteFile", err)
		return err
	}

	return nil
}
