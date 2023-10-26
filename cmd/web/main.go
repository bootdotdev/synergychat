package main

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed public/*
var content embed.FS

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

	// Setting the content type based on the file extension
	if ext := http.DetectContentType(data); len(ext) > 0 {
		w.Header().Set("Content-Type", ext)
	}

	w.Write(data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Serving static SynergyChat front-end on :3000")
	http.ListenAndServe(":3000", nil)
}
