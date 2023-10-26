package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Book struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Subjects []string `json:"subjects"`
	Authors  []struct {
		Name string `json:"name"`
	} `json:"authors"`
	Languages     []string          `json:"languages"`
	MediaType     string            `json:"media_type"`
	Formats       map[string]string `json:"formats"`
	DownloadCount int               `json:"download_count"`
}

type Match struct {
	BookTitle string
	Count     int
}

const timeBetweenBooks = time.Second

func (cfg apiConfig) worker() {
	nextURL := cfg.baseURL
	for nextURL != "" {
		books, next, err := fetchBooks(nextURL)
		if err != nil {
			log.Println("Error fetching books:", err)
			return
		}

		for _, book := range books {
			count := checkBookForKeywords(book, cfg.keywords)
			log.Printf("Found %v matches in %v", count, book.Title)
			if count > 0 {
				err = cfg.db.saveCount(book.Title, count)
				if err != nil {
					log.Printf("Error saving count for %v: %v", book.Title, err)
				}
			}
			time.Sleep(timeBetweenBooks)
		}
		nextURL = next
	}
}

func fetchBooks(nextURL string) ([]Book, string, error) {
	log.Println("Fetching books from", nextURL)
	resp, err := http.Get(nextURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		Results []Book `json:"results"`
		Next    string `json:"next"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, "", err
	}

	return result.Results, result.Next, nil
}

func checkBookForKeywords(book Book, keywords []string) int {
	log.Println("Checking book", book.Title)
	textUrl, ok := book.Formats["text/plain"]
	if !ok {
		return 0
	}

	resp, err := http.Get(textUrl)
	if err != nil {
		log.Println("Error fetching book text:", err)
		return 0
	}
	defer resp.Body.Close()

	text, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading book text:", err)
		return 0
	}

	count := 0
	lines := strings.Split(string(text), "\n")
	for _, line := range lines {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				count++
			}
		}
	}
	return count
}
