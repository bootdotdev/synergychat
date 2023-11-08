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

const timeBetweenBooks = time.Second * 30
const baseURLFormat = "http://gutendex.com/books/?page=%d"

func (cfg apiConfig) worker() {
	nextURL := cfg.baseURL
	for nextURL != "" {
		books, next, err := fetchBooks(nextURL)
		if err != nil {
			log.Println("Error fetching books:", err)
			return
		}

		for _, book := range books {
			counts := checkBookForKeywords(book, cfg.keywords)
			for keyword, count := range counts {
				err = cfg.db.saveCount(keyword, book.Title, count)
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

func checkBookForKeywords(book Book, keywords []string) map[string]int {
	log.Println("Checking book", book.Title)
	textUrl, ok := book.Formats["text/plain"]
	if !ok {
		return nil
	}

	resp, err := http.Get(textUrl)
	if err != nil {
		log.Println("Error fetching book text:", err)
		return nil
	}
	defer resp.Body.Close()

	text, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading book text:", err)
		return nil
	}

	counts := make(map[string]int)
	lines := strings.Split(string(text), "\n")
	for _, line := range lines {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				if _, ok := counts[keyword]; !ok {
					counts[keyword] = 0
				}
				counts[keyword]++
			}
		}
	}
	return counts
}
