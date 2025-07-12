package constData

import (
	"errors"
	"fmt"
)

// ----- Errors -----
var (
	// ErrConvString is returned when a string cannot be converted
	// (e.g. "&&##" has no Morse equivalent).
	ErrConvString = errors.New("unable to convert string")

	// ErrEmptyInput is returned when an empty argument is supplied.
	// For example it is not possible to code/decode empty string to Morse.
	ErrEmptyInput = errors.New("input is empty")

	// ErrFileCreate is returned when it is not possible to create file.
	ErrFileCreate = errors.New("file creation error")
)

// ----- Error messages formated -----
var (
	ErrEmptyFmt         = "%w: conversion result is empty\n"
	ErrEmptyMorseString = "%w: not possible to code/decode empty string to Morse\n"
	ErrWrongSymbol      = "%w: symbol %q not allowed\n"
	ErrAmbivalentInput  = "%w: unable to determine whether the input is Morse code or text\n"
)

// ----- Error messages constant -----
var (
	ErrIndexHTMLMissing = "Missing 'index.html' in work directory"
	ErrIndexHTMLParsing = "'index.html' is broken"
	ErrLogCreate        = "not possible to create log file"
	ErrServCreation     = "error during server creation"
	ErrFileNotFound     = "file not found"
	ErrFileCreateMsg    = "file create error"
	ErrBadFileContent   = "bad file content"
	ErrFileWrite        = "file write error"
	ErrConvStringMsg    = "unable to convert string"
)

// ----- Server parameters -----
const (
	DefaultPort         = 8080
	DefaultReadTimeout  = 5  // seconds
	DefaultWriteTimeout = 10 // seconds
	DefaultIdleTimeout  = 15 // seconds
)

// ----- Log parameters -----
var (
	LogFileName         = "morseServer.log"
	LogPrefix           = "LOG: "
	LogSuccServCreation = fmt.Sprintf("server started! address: localhost:%d", DefaultPort)
)

// ----- Morse code parameters -----
// allowedCharactersText returns a string containing all runes considered
// valid for text input (Cyrillic letters, digits and punctuation
// supported by pkg/morse).
// Since this function is
func allowedCharactersText() string {

	// adding symbols compatible with morse.go
	var ch = []rune{'"', '\'', '(', ')', ',', '-', '.', '/', ':', '?', ' '}

	// adding symbols from Russian alphabet
	for i := 1040; i <= 1103; i++ {
		ch = append(ch, rune(i))
	}

	// adding digits
	for i := 48; i <= 57; i++ {
		ch = append(ch, rune(i))
	}

	return string(ch)
}

var (
	AllwText  = allowedCharactersText()
	AllwMorse = ".- "
)

// ----- Handlers data ------
var IndexHTML = "index.html"
