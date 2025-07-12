package logger

import (
	"io"
	"log"
	"os"

	con "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/constData"
)

// lg хранит ссылку на общий логгер.
// По умолчанию — стандартный log.Default().
var Lg = log.Default()

// SetLogger позволяет приложению заменить логгер.
func SetLogger(l *log.Logger) {
	if l != nil {
		Lg = l
	}
}

func CreateLogger() (*log.Logger, *os.File) {

	file, err := os.OpenFile(con.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%v: %v", con.ErrLogCreate, con.ErrFileCreate)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := log.New(multiWriter, con.LogPrefix, log.Ldate|log.Ltime)

	return logger, file
}
