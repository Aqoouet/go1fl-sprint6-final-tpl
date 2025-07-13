// Package service provides function for automatic detection and
// conversion between text and Morse code.

package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// ----- Errors -----
var (
	// ErrConvString is returned when a string cannot be converted
	// (e.g. "&&##" has no Morse equivalent).
	ErrConvString = errors.New("unable to convert string")

	// ErrEmptyInput is returned when an empty argument is supplied.
	// For example it is not possible to code/decode empty string to Morse.
	ErrEmptyInput = errors.New("input is empty")
)

// ----- Error messages formated -----
var (
	ErrEmptyFmt         = "%w: conversion result is empty\n"
	ErrEmptyMorseString = "%w: not possible to code/decode empty string to Morse\n"
	ErrWrongSymbol      = "%w: symbol %q not allowed\n"
	ErrAmbivalentInput  = "%w: unable to determine whether the input is Morse code or text\n"
)

// ConvertDetectMorse automatically detects whether input is Morse code or text
// and converts it to the opposite representation
// Returns error if input contains invalid symbols or cannot be converted
func ConvertDetectMorse(s string) (string, error) {

	if s == "" {
		return "", fmt.Errorf(ErrEmptyMorseString, ErrEmptyInput)
	}

	isMorse := true
	isText := true

	wrongSymbol := ""

	for _, c := range s {

		if !strings.ContainsRune(morse.AllwMorse, c) {
			isMorse = false
		}
		if !strings.ContainsRune(morse.AllwText, c) {
			isText = false
			wrongSymbol = string(c)
			break
		}
	}

	if wrongSymbol != "" {
		return "", fmt.Errorf(ErrWrongSymbol, ErrConvString, wrongSymbol)
	}

	if isMorse {
		if t := morse.ToText(s); t == "" {
			return "", fmt.Errorf(ErrEmptyFmt, ErrConvString)
		} else {
			return t, nil
		}
	}

	if isText {
		if t := morse.ToMorse(s); t == "" {
			return "", fmt.Errorf(ErrEmptyFmt, ErrConvString)
		} else {
			return t, nil
		}
	}

	return "", fmt.Errorf(ErrAmbivalentInput, ErrConvString)
}
