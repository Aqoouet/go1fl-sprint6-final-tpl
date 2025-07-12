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

func indexHTMLParse() []byte {

	indexBytes, err := os.ReadFile(con.IndexHTML)
	if err != nil {
		logger.Lg.Fatalf("%v: %v", con.ErrIndexHTMLParsing, err)
	}

	return indexBytes

}

var indexBytes = indexHTMLParse()

func GetMain(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(indexBytes)

}

func PostUpload(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(0); err != nil {
		logger.Lg.Printf("%v: %v\n", con.ErrIndexHTMLParsing, err)
		http.Error(w, con.ErrIndexHTMLParsing, http.StatusBadRequest)
		return
	}

	var fi multipart.File
	var err error

	for k := range r.MultipartForm.File {
		fi, _, err = r.FormFile(k)
		if err == nil {
			break
		}
	}

	if err != nil {
		logger.Lg.Printf("%v: %v\n", con.ErrFileNotFound, err)
		http.Error(w, con.ErrFileNotFound, http.StatusBadRequest)
		return
	}

	fileData, err := io.ReadAll(fi)
	fi.Close()
	if err != nil {
		http.Error(w, con.ErrBadFileContent, http.StatusBadRequest)
		logger.Lg.Printf("%v: %v\n", con.ErrBadFileContent, err)
		return
	}

	converted, err := service.ConvertDetectMorse(string(fileData))
	if err != nil {
		http.Error(w, con.ErrConvStringMsg, http.StatusBadRequest)
		logger.Lg.Printf("%v: %v\n", con.ErrConvStringMsg, err)
	}

	outputFileName := fmt.Sprintf("conversion result %s", time.Now().UTC().Format("2006-01-02 15:04:05"))
	d, err := os.Create(outputFileName)
	if err != nil {
		http.Error(w, con.ErrFileCreateMsg, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", con.ErrFileCreateMsg, err)
	}
	defer d.Close()

	_, err = io.WriteString(d, converted)
	if err != nil {
		http.Error(w, con.ErrFileWrite, http.StatusInternalServerError)
		logger.Lg.Printf("%v: %v\n", con.ErrFileWrite, err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "saved as %s\n", outputFileName)

	logger.Lg.Printf("uploaded file converted -> %s", outputFileName)

}
