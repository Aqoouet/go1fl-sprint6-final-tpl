// Package logger provides logging functionality for the application.
// It includes logger creation, configuration, and global logger management
// with support for both console and file output.
package logger

import (
	"io"
	"log"
	"os"
)

// ----- Log parameters -----
var (
	LogFileName   = "morseServer.log"
	LogPrefix     = "LOG: "
	ErrLogCreate  = "not possible to create log file"
	ErrFileCreate = "file creation error"
)

// Lg stores reference to the global logger
// By default uses log.Default()
var Lg = log.Default()

// SetLogger allows the application to replace the logger
// Only replaces if the provided logger is not nil
func SetLogger(l *log.Logger) {
	if l != nil {
		Lg = l
	}
}

// CreateLogger creates a new logger that writes to both stdout and log file
// Returns the logger instance and the log file handle
// Log file is created with append mode and 0666 permissions
func CreateLogger() (*log.Logger, *os.File) {

	file, err := os.OpenFile(LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%v: %v", ErrLogCreate, ErrFileCreate)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := log.New(multiWriter, LogPrefix, log.Ldate|log.Ltime)

	return logger, file
}
