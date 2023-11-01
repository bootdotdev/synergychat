package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	keywords []string
	baseURL  string
	db       DB
}

func main() {
	// "https://gutendex.com/books"
	baseURL := os.Getenv("TO_CRAWL_URL")
	if baseURL == "" {
		log.Fatal("No TO_CRAWL_URL found in environment")
	}
	// "love,hate"
	keywordsString := os.Getenv("CRAWLER_KEYWORDS")
	if keywordsString == "" {
		log.Fatal("No CRAWLER_KEYWORDS found in environment")
	}
	keywords := strings.Split(keywordsString, ",")

	// 5000
	port := os.Getenv("CRAWLER_PORT")
	if port == "" {
		log.Fatal("No CRAWLER_PORT found in environment")
	}

	apiCfg := apiConfig{
		keywords: keywords,
		baseURL:  baseURL,
	}

	// optional
	// e.g. "./crawler-db"
	crawlerDBPath := os.Getenv("CRAWLER_DB_PATH")
	if crawlerDBPath == "" {
		var mem *Memory
		apiCfg.db = mem
	} else {
		apiCfg.db = &Disk{
			crawlerDBPath: crawlerDBPath,
		}
	}

	err := apiCfg.db.init()
	if err != nil {
		log.Fatal("Couldn't initialize database: ", err)
	}
	go apiCfg.worker()

	router := chi.NewRouter()
	router.Get("/healthz", handlerReadiness)
	router.Get("/stats", apiCfg.handlerGetAllMatches)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg apiConfig) handlerGetAllMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := cfg.db.getCounts()
	if err != nil {
		respondWithJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	if r.URL.Query().Get("keywords") != "" {
		matches = filterKeywords(r.URL.Query().Get("keywords"), matches)
	}
	if r.URL.Query().Get("title") != "" {
		matches = filterTitles(r.URL.Query().Get("title"), matches)
	}

	matchesSlice := matchesMapToSlice(matches)
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
