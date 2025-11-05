package main

import (
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"

	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "MORSE_CONVERTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := handlers.Init(); err != nil {
		logger.Fatalf("Failed to initialize handlers: %v", err)
	}
	srv := server.New(logger)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
