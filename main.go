// http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
package main

import (
	"bytes"
	"encoding/json"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/davidwalter0/go-cfg"
	"golang.org/x/net/http2"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type App struct {
	Cert  string
	Key   string
	Host  string
	Port  string
	HTTPS bool
}

var app App

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	QrGenerator(w, r)
}

func run() {
	var jsonText []byte
	var err error
	if err = cfg.Parse(&app); err != nil {
		log.Fatalf("%v\n", err)
	}

	jsonText, _ = json.MarshalIndent(&app, "", "  ")
	protocol := "http://"
	if app.HTTPS {
		protocol = "https://"
	}
	log.Printf("\nEnvironment configuration\n%v\n", string(jsonText))

	log.Println("Answering requests on " + protocol + app.Host + ":" + app.Port)
	if app.HTTPS {
		handler := MyHandler{}
		server := http.Server{
			Addr:    app.Host + ":" + app.Port,
			Handler: &handler,
		}
		http2.ConfigureServer(&server, &http2.Server{})
		err = server.ListenAndServeTLS(app.Cert, app.Key)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		http.HandleFunc("/", QrGenerator)
		err = http.ListenAndServe(app.Host+":"+app.Port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func main() {
	run()
}

func QrGenerator(w http.ResponseWriter, r *http.Request) {
	log.Println("http://"+app.Host+":"+app.Port, r.Body, r)
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
