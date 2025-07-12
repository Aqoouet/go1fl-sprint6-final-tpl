// Package service provides function for automatic detection and
// conversion between text and Morse code.

package service

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"

	сon "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/constData"
)

// convertDetectMorse detects whether s is Morse code or text and
// converts it to the opposite representation. It returns an error if
// the input contains invalid symbols or cannot be converted.
func ConvertDetectMorse(s string) (string, error) {

	if s == "" {
		return "", fmt.Errorf(сon.ErrEmptyMorseString, сon.ErrEmptyInput)
	}

	isMorse := true
	isText := true

	wrongSymbol := ""

	for _, c := range s {

		if !strings.ContainsRune(сon.AllwMorse, c) {
			isMorse = false
		}
		if !strings.ContainsRune(сon.AllwText, c) {
			isText = false
			wrongSymbol = string(c)
			break
		}
	}

	if wrongSymbol != "" {
		return "", fmt.Errorf(сon.ErrWrongSymbol, сon.ErrConvString, wrongSymbol)
	}

	if isMorse {
		if t := morse.ToText(s); t == "" {
			return "", fmt.Errorf(сon.ErrEmptyFmt, сon.ErrConvString)
		} else {
			return t, nil
		}
	}

	if isText {
		if t := morse.ToMorse(s); t == "" {
			return "", fmt.Errorf(сon.ErrEmptyFmt, сon.ErrConvString)
		} else {
			return t, nil
		}
	}

	return "", fmt.Errorf(сon.ErrAmbivalentInput, сon.ErrConvString)
}
