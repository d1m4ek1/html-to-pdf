package logging

import (
	"fmt"
	"html-to-pdf/pkg/seterror"
	"os"
	"time"
)

type Logging struct {
	Filename        string
	Date            string
	OperationTime   string
	OperationMemory int64
	WarningMessage  string
}

type CreateLogging interface {
	wrapLogComplete() string

	wrapLogWarning() string

	SetLogComplete() error

	SetLogWarning(warningMessage string) error
}

const (
	// pathToLogCompleteOperations Путь к лог файлу с завершенными операциями
	pathToLogCompleteOperations string = "logs/completedOperations.log"

	// pathToLogWarnings Путь к лог файлу с предупреждающими логами
	pathToLogWarnings string = "logs/warnings.log"
)

// createDate Создает текущую дату
func createDate() string {
	var dateTime time.Time = time.Now()
	var date string = fmt.Sprintf("%d.%s.%d", dateTime.Day(), dateTime.Month().String(), dateTime.Year())
	return date
}

// wrapLogComplete Пакует шаблон для лога
func (l *Logging) wrapLogComplete() string {
	var formatString string = fmt.Sprintf(
		`
Complete operation:
	File name: %s,
	Date of operation: %s,
	Time of operation: %s,
	Memory of operation: %d;
`, l.Filename, l.Date, l.OperationTime, l.OperationMemory)

	return formatString
}

// wrapLogWarning Пакует шаблон для лога
func (l *Logging) wrapLogWarning() string {
	var formatString string = fmt.Sprintf(
		`
Warning:
	Message: %s,
	Date: %s`, l.WarningMessage, l.Date)

	return formatString
}

// SetLogComplete Записывает лог о завершении операции в файл completedOperations.log
func (l *Logging) SetLogComplete() error {
	file, err := os.OpenFile(pathToLogCompleteOperations, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		seterror.SetAppError("os.Open", err)
		return err
	}

	if _, err := file.WriteString(l.wrapLogComplete()); err != nil {
		seterror.SetAppError("file.WriteString", err)
		return err
	}

	return nil
}

// SetLogWarning Записывает предупреждающий лог в файл warnings.log
func (l *Logging) SetLogWarning(warningMessage string) error {
	file, err := os.OpenFile(pathToLogWarnings, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		seterror.SetAppError("os.Open", err)
		return err
	}

	l.WarningMessage = warningMessage

	if _, err := file.WriteString(l.wrapLogWarning()); err != nil {
		seterror.SetAppError("file.WriteString", err)
		return err
	}

	return nil
}

// NewLogging Создает структуру Logging
func NewLogging(fileName, operationTime string, operationMemory int64) *Logging {
	return &Logging{
		Filename:        fileName,
		Date:            createDate(),
		OperationTime:   operationTime,
		OperationMemory: operationMemory,
	}
}
