package createpdf

import (
	"archive/zip"
	"html-to-pdf/pkg/seterror"
	"io"
	"os"
)

// loadPictures Загружает и сохраняет изображения в временном хранилище
func loadPictures(file *zip.File, pathToTemporaryStorage string) error {
	destinationFile, err := os.OpenFile(pathToTemporaryStorage+"/"+file.FileInfo().Name(),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
	if err != nil {
		seterror.SetAppError("os.OpenFile", err)
		return err
	}

	contentReadCloser, err := file.Open()
	if err != nil {
		seterror.SetAppError("file.Open", err)
		return err
	}

	if _, err := io.Copy(destinationFile, contentReadCloser); err != nil {
		seterror.SetAppError("io.Copy", err)
		return err
	}

	return nil
}
