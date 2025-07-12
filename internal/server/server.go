// Package server provides HTTP server creation and configuration functionality.
// It includes server setup with routing, timeout configuration, and logger integration
// for the Morse code conversion service.
package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/logger"
	"github.com/go-chi/chi/v5"

	con "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/constData"
)

// morseServer represents the HTTP server with logger and server instance
type morseServer struct {
	Log  *log.Logger
	Serv *http.Server
}

// newRouter creates and configures the HTTP router with registered handlers
// Sets up logger and registers GET "/" and POST "/upload" endpoints
func newRouter(appLog *log.Logger) http.Handler {

	logger.SetLogger(appLog)

	r := chi.NewRouter()

	r.Get("/", handlers.GetMain)
	r.Post("/upload", handlers.PostUpload)

	return r
}

// CreateRouter creates and starts a new HTTP server with configured settings
// Returns a morseServer instance with logger and server components
// Server runs on port 8080 with configured timeouts and error logging
func CreateRouter(l *log.Logger) *morseServer {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", con.DefaultPort),
		Handler:      newRouter(l),
		ReadTimeout:  con.DefaultReadTimeout * time.Second,
		WriteTimeout: con.DefaultWriteTimeout * time.Second,
		IdleTimeout:  con.DefaultIdleTimeout * time.Second,
		ErrorLog:     l,
	}

	// Start server in a goroutine to avoid blocking the main thread
	// This allows the function to return immediately while the server runs in background
	// useful in tests
	go func() {

		l.Println(con.LogSuccServCreation)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Printf("%v: %v", con.ErrServCreation, err.Error())
		}

	}()

	return &morseServer{l, server}
}
