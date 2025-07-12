package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	сon "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/constData"
)

func runConvertDetectMorseTest(input string, output string, hasError bool, typeError error, t *testing.T) {
	result, err := ConvertDetectMorse(input)
	if hasError {
		require.Error(t, err)
		assert.ErrorIs(t, err, typeError)
	} else {
		require.NoError(t, err)
		assert.Equal(t, output, result)
	}
}

type convertDetectMorseTestItem struct {
	name      string
	input     string
	output    string
	hasError  bool
	typeError error
}

var tests = []convertDetectMorseTestItem{
	{
		name:      "морзе",
		input:     ".- -... .-- --.",
		output:    "АБВГ",
		hasError:  false,
		typeError: nil,
	},
	{
		name:      "обычный текст",
		input:     "ПРИВЕТ МИР",
		output:    ".--. .-. .. .-- . - -- .. .-.",
		hasError:  false,
		typeError: nil,
	},
	{
		name:      "недопустимые символы",
		input:     "строка содержит символы % ! @",
		output:    "",
		hasError:  true,
		typeError: сon.ErrConvString,
	},
	{
		name:      "пустая строка",
		input:     "",
		output:    "",
		hasError:  true,
		typeError: сon.ErrEmptyInput,
	},
	{
		name:      "текст содержащий строчные буквы с цифрами",
		input:     "Москва 1980",
		output:    "-- --- ... -.- .-- .- .---- ----. ---.. -----",
		hasError:  false,
		typeError: nil,
	},
	{
		name:      "морзе с цифрами",
		input:     ".-.. --- -. -.. --- -.   ..--- ----- .---- ..---",
		output:    "ЛОНДОН 2012",
		hasError:  false,
		typeError: nil,
	},
	{
		name:      "неструктурированный код Морзе",
		input:     ".-..----.-..----...--------.----..---",
		output:    "",
		hasError:  true,
		typeError: сon.ErrConvString,
	},
	{
		name:      "одна точка",
		input:     ".",
		output:    "Е",
		hasError:  false,
		typeError: nil,
	},
	{
		name:      "текст с латиницей",
		input:     "Hello World",
		output:    "",
		hasError:  true,
		typeError: сon.ErrConvString,
	},
}

func TestConvertDetectMorse(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runConvertDetectMorseTest(tt.input, tt.output, tt.hasError, tt.typeError, t)
		})
	}
}
