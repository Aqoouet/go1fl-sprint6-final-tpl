// Package handlers provides HTTP request handlers for the Morse code conversion service.
// It includes handlers for serving the main HTML page and processing file uploads
// with automatic Morse code detection and conversion.
package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/logger"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"

	con "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/constData"
)

// ----- Index html -----

// indexHTMLParse reads and returns the content of index.html file
// Returns the HTML content as bytes or terminates the program on error
func indexHTMLParse() []byte {

	indexBytes, err := os.ReadFile(con.IndexHTML)
	if err != nil {
		logger.Lg.Fatalf("%v: %v", con.ErrIndexHTMLParsing, err)
	}

	return indexBytes

}

// getIndexBytes returns the content of index.html file
// index.html is parsed only when it is necessary
// useful tests: it is possible to replace index.html by another file
func getIndexBytes() []byte {
	return indexHTMLParse()
}

// GetMain handles the root endpoint "/" and returns the main HTML page
// Sets appropriate content-type header and returns the index.html content
func GetMain(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(getIndexBytes())

}

// PostUpload handles file upload and Morse code conversion
// It parses multipart form data, reads the uploaded file,
// converts the content using automatic Morse detection,
// saves the result to a local file, and returns the conversion details
func PostUpload(w http.ResponseWriter, r *http.Request) {

	// Parse multipart form data from the request
	// This extracts the uploaded file from the HTML form
	if err := r.ParseMultipartForm(0); err != nil {
		logger.Lg.Printf("%v: %v\n", con.ErrIndexHTMLParsing, err)
		http.Error(w, con.ErrIndexHTMLParsing, http.StatusBadRequest)
		return
	}

	var fi multipart.File
	var err error

	// Iterate through all uploaded files in the form
	// We take the first available file (typically there's only one)
	for k := range r.MultipartForm.File {
		fi, _, err = r.FormFile(k)
		if err == nil {
			break
		}
	}

	// Check if file was found and successfully opened
	if err != nil || fi == nil {
		logger.Lg.Printf("%v: %v\n", con.ErrFileNotFound, err)
		http.Error(w, con.ErrFileNotFound, http.StatusBadRequest)
		return
	}

	// Read all data from the uploaded file
	fileData, err := io.ReadAll(fi)
	fi.Close() // Always close the file to free resources
	if err != nil {
		http.Error(w, con.ErrBadFileContent, http.StatusBadRequest)
		logger.Lg.Printf("%v: %v\n", con.ErrBadFileContent, err)
		return
	}

	// Convert file content to string and detect/convert Morse code
	// The service automatically determines if input is text or Morse code
	initial_string := string(fileData)
	converted, err := service.ConvertDetectMorse(initial_string)
	if err != nil {
		http.Error(w, con.ErrConvStringMsg, http.StatusBadRequest)
		logger.Lg.Printf("%v: %v\n", con.ErrConvStringMsg, err)
		return
	}

	// Generate unique filename for the conversion result
	// Format: conversion_result_YYYY-MM-DD_HH:MM:SS
	outputFileName := fmt.Sprintf("conversion_result_%s", time.Now().UTC().Format("2006-01-02_15:04:05"))

	// Create the output file to save conversion results
	d, err := os.Create(outputFileName)
	if err != nil {
		http.Error(w, con.ErrFileCreateMsg, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", con.ErrFileCreateMsg, err)
		return
	}
	defer d.Close() // Ensure file is closed even if errors occur

	// Write the converted content to the output file
	_, err = io.WriteString(d, converted)
	if err != nil {
		http.Error(w, con.ErrFileWrite, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", con.ErrFileWrite, err)
		return
	}

	// Set HTTP status to 201 (Created) to indicate successful file creation
	w.WriteHeader(http.StatusCreated)

	// Return detailed information about the conversion
	// This includes both the original input and the converted output
	fmt.Fprintf(w, "initial string %q was converted to %q\n", initial_string, converted)
	fmt.Fprintf(w, "results of conversion were saved in file %q\n", outputFileName)

	// Log successful conversion for debugging/monitoring
	logger.Lg.Printf("uploaded file converted -> %s", outputFileName)

}
