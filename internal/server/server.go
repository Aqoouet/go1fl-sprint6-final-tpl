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

type morseServer struct {
	Log  *log.Logger
	Serv *http.Server
}

func newRouter(appLog *log.Logger) http.Handler {

	logger.SetLogger(appLog)

	r := chi.NewRouter()

	r.Get("/", handlers.GetMain)
	r.Post("/upload", handlers.PostUpload)

	return r
}

func CreateRouter(l *log.Logger) *morseServer {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", con.DefaultPort),
		Handler:      newRouter(l),
		ReadTimeout:  con.DefaultReadTimeout * time.Second,
		WriteTimeout: con.DefaultWriteTimeout * time.Second,
		IdleTimeout:  con.DefaultIdleTimeout * time.Second,
		ErrorLog:     l,
	}

	go func() {

		l.Println(con.LogSuccServCreation)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Printf("%v: %v", con.ErrServCreation, err.Error())
		}

	}()

	return &morseServer{l, server}
}
