package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type apiConfig struct {
	db         DB
	crawlerURL string
}

func main() {
	// 5000
	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("No API_PORT found in environment")
	}

	apiCfg := apiConfig{
		crawlerURL: os.Getenv("CRAWLER_BASE_URL"),
	}
	apiDBFilePath := os.Getenv("API_DB_FILEPATH")
	if apiDBFilePath == "" {
		apiCfg.db = &Memory{}
	} else {
		apiCfg.db = &Disk{
			apiDBFilePath: apiDBFilePath,
		}
	}
	err := apiCfg.db.init()
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/healthz", handlerReadiness)
	router.Post("/messages", apiCfg.handlerCreateMessage)
	router.Get("/messages", apiCfg.handlerGetMessages)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (cfg apiConfig) handlerCreateMessage(w http.ResponseWriter, r *http.Request) {
	reqBody := Message{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithJSON(w, 400, map[string]string{
			"error": fmt.Sprintf("Error decoding JSON: %s", err),
		})
		return
	}

	if containsProfanity(reqBody.Text) {
		log.Fatal("Profanity detected!!!")
	}

	err = cfg.db.createMessage(reqBody.AuthorUsername, reqBody.Text)
	if err != nil {
		respondWithJSON(w, 400, map[string]string{
			"error": fmt.Sprintf("Error creating message: %s", err),
		})
		return
	}

	err = cfg.handleSlashCommand(reqBody.Text)
	if err != nil {
		respondWithJSON(w, 400, map[string]string{
			"error": fmt.Sprintf("Error handling slash command: %s", err),
		})
		return
	}

	messages, err := cfg.db.getMessages()
	if err != nil {
		respondWithJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respondWithJSON(w, http.StatusCreated, messages)
}

func (cfg apiConfig) handlerGetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := cfg.db.getMessages()
	if err != nil {
		respondWithJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respondWithJSON(w, 200, messages)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func containsProfanity(message string) bool {
	lowered := strings.ToLower(message)
	if strings.Contains(lowered, "darn") {
		return true
	}
	if strings.Contains(lowered, "heck") {
		return true
	}
	if strings.Contains(lowered, "fetch") {
		return true
	}
	return false
}
