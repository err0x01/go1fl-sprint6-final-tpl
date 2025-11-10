package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "MORSE_CONVERTER: ", log.Ldate|log.Ltime|log.Lshortfile)
	srv := server.New(logger)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
