// http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
package main

import (
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var host = os.Getenv("QRCODE_GENERATOR_HOST")
var port = os.Getenv("QRCODE_GENERATOR_PORT")

func main() {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	log.Println("Answering requests on on http://" + host + ":" + port)
	http.HandleFunc("/", QrGenerator)

	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func QrGenerator(w http.ResponseWriter, r *http.Request) {
	log.Println("http://"+host+":"+port, r.Body, r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	data := r.URL.Query().Get("data")
	if data == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	s, err := url.QueryUnescape(data)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	code, err := qr.Encode(s, qr.L, qr.Auto)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	size := r.URL.Query().Get("size")
	if size == "" {
		size = "250"
	}
	intsize, err := strconv.Atoi(size)
	if err != nil {
		intsize = 250
	}

	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, intsize, intsize)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, code); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
