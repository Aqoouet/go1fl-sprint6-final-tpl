// Package main is the entry point for the Morse code conversion service.
// It initializes the logger, creates and starts the HTTP server,
// and manages the application lifecycle.
package main

import (
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/logger"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	l, f := logger.CreateLogger()
	defer f.Close()
	_ = server.CreateRouter(l)

	select {}
}
