package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
)

//go:embed public/*
var content embed.FS

var apiURL = ""

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	data, err := content.ReadFile("public" + path)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}

	// add a string to the javascript file before serving it
	if path == "/app.js" {
		data = append([]byte(fmt.Sprintf("const apiUrl = '%s';\n\n", apiURL)), data...)
	}

	// Set headers to disable caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0
	w.Header().Set("Expires", "0")

	// Setting the content type based on the file extension
	if ext := http.DetectContentType(data); len(ext) > 0 {
		w.Header().Set("Content-Type", ext)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	}

	w.Write(data)
}

func main() {
	// 5007
	port := os.Getenv("WEB_PORT")
	if port == "" {
		log.Fatal("No WEB_PORT found in environment")
	}

	apiURL = os.Getenv("API_URL")
	if port == "" {
		log.Fatal("No API_URL found in environment")
	}

	http.HandleFunc("/", handler)
	fmt.Println("Serving static SynergyChat front-end on :" + port)
	http.ListenAndServe(":"+port, nil)
}
