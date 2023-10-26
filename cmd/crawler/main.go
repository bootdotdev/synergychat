package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	matches   map[string]Match
	matchesMu *sync.Mutex
	keywords  []string
	baseURL   string
}

func main() {
	// e.g. "https://gutendex.com"
	baseURL := os.Getenv("CRAWLER_BASE_URL")
	if baseURL == "" {
		log.Fatal("No CRAWLER_BASE_URL found in environment")
	}
	// e.g. "love,hate"
	keywordsString := os.Getenv("CRAWLER_KEYWORDS")
	if keywordsString == "" {
		log.Fatal("No CRAWLER_KEYWORDS found in environment")
	}
	keywords := strings.Split(keywordsString, ",")

	port := os.Getenv("CRAWLER_PORT")
	if port == "" {
		log.Fatal("No CRAWLER_PORT found in environment")
	}

	apiCfg := apiConfig{
		matches:   map[string]Match{},
		matchesMu: &sync.Mutex{},
		keywords:  keywords,
		baseURL:   baseURL,
	}
	go apiCfg.worker()

	router := chi.NewRouter()
	router.Get("/healthz", handlerReadiness)
	router.Get("/matches", apiCfg.handlerGetMatches)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg apiConfig) handlerGetMatches(w http.ResponseWriter, r *http.Request) {
	cfg.matchesMu.Lock()
	defer cfg.matchesMu.Unlock()

	matchesSlice := []Match{}
	for _, match := range cfg.matches {
		matchesSlice = append(matchesSlice, match)
	}
	slices.SortFunc(matchesSlice, func(a, b Match) int {
		return a.Count - b.Count
	})

	respondWithJSON(w, 200, matchesSlice)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
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
