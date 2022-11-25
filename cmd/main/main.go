package main

import (
	"html-to-pdf/internal/handlers"
	"html-to-pdf/pkg/seterror"
	"log"
)

func main() {
	log.Println("Start service")
	if err := handlers.InitGin(); err != nil {
		seterror.SetAppError("handlers.InitGin", err)
		return
	}
}
