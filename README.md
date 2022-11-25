# html-to-pdf service

Версия golang 1.19

## Дополнительно установить к проекту
sudo apt install wkhtmltopdf

go get -u github.com/SebastiaanKlippert/go-wkhtmltopdf

go get github.com/gin-gonic/gin


## API

POST /api/upload - загружает zip-файл и конвертирует файлы в pdf файл

## Запуск
go run ./cmd/web/main.go, либо make start