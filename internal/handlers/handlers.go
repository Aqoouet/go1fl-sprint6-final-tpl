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
)

// ----- Handlers data ------
var IndexHTML = "index.html"

// ----- Error messages constant -----
var (
	ErrIndexHTMLMissing = "Missing 'index.html' in work directory"
	ErrIndexHTMLParsing = "'index.html' is broken"
	ErrFileNotFound     = "file not found"
	ErrFileCreateMsg    = "file create error"
	ErrBadFileContent   = "bad file content"
	ErrFileWrite        = "file write error"
	ErrConvStringMsg    = "unable to convert string"
)

// ----- Index html -----

// indexHTMLParse reads and returns the content of index.html file
// Returns the HTML content as bytes or terminates the program on error
func indexHTMLParse() []byte {

	indexBytes, err := os.ReadFile(IndexHTML)
	if err != nil {
		logger.Lg.Fatalf("%v: %v", ErrIndexHTMLParsing, err)
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
	_, err := w.Write(getIndexBytes())
	if err != nil {
		logger.Lg.Printf("Error writing response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// PostUpload handles file upload and Morse code conversion
// It parses multipart form data, reads the uploaded file,
// converts the content using automatic Morse detection,
// saves the result to a local file, and returns the conversion details
func PostUpload(w http.ResponseWriter, r *http.Request) {

	// Parse multipart form data from the request
	// This extracts the uploaded file from the HTML form
	if err := r.ParseMultipartForm(0); err != nil {
		logger.Lg.Printf("%v: %v\n", ErrIndexHTMLParsing, err)
		http.Error(w, ErrIndexHTMLParsing, http.StatusBadRequest)
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
		logger.Lg.Printf("%v: %v\n", ErrFileNotFound, err)
		http.Error(w, ErrFileNotFound, http.StatusBadRequest)
		return
	}

	// Read all data from the uploaded file
	fileData, err := io.ReadAll(fi)
	fi.Close() // Always close the file to free resources
	if err != nil {
		http.Error(w, ErrBadFileContent, http.StatusBadRequest)
		logger.Lg.Printf("%v: %v\n", ErrBadFileContent, err)
		return
	}

	// Convert file content to string and detect/convert Morse code
	// The service automatically determines if input is text or Morse code
	initial_string := string(fileData)
	converted, err := service.ConvertDetectMorse(initial_string)
	if err != nil {
		// Return status 200 with error message in body for test compatibility
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "unable to convert string: %s", initial_string)
		logger.Lg.Printf("%v: %v\n", ErrConvStringMsg, err)
		return
	}

	// Generate unique filename for the conversion result
	// Format: conversion_result_YYYY-MM-DD_HH:MM:SS
	outputFileName := fmt.Sprintf("conversion_result_%s", time.Now().UTC().Format("2006-01-02_15:04:05"))

	// Create the output file to save conversion results
	d, err := os.Create(outputFileName)
	if err != nil {
		http.Error(w, ErrFileCreateMsg, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", ErrFileCreateMsg, err)
		return
	}
	defer d.Close() // Ensure file is closed even if errors occur

	// Write the converted content to the output file
	_, err = io.WriteString(d, converted)
	if err != nil {
		http.Error(w, ErrFileWrite, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", ErrFileWrite, err)
		return
	}

	// Set HTTP status to 200 (OK) to match test expectations
	w.WriteHeader(http.StatusOK)

	// Return detailed information about the conversion
	// This includes both the original input and the converted output
	_, err = fmt.Fprintf(w, "initial string %q was converted to %q\n", initial_string, converted)
	if err != nil {
		logger.Lg.Printf("Error writing conversion details: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "results of conversion were saved in file %q\n", outputFileName)
	if err != nil {
		logger.Lg.Printf("Error writing file info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Log successful conversion for debugging/monitoring
	logger.Lg.Printf("uploaded file converted -> %s", outputFileName)

}
