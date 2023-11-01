package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const crawlerBotAuthorUsername = "crawler-bot"

type Match struct {
	Keyword   string
	BookTitle string
	Count     int
}

func (cfg apiConfig) handleSlashCommand(msg string) error {
	if !strings.HasPrefix(msg, "/stats") {
		return nil
	}
	if cfg.crawlerURL == "" {
		err := cfg.db.createMessage(crawlerBotAuthorUsername, "Crawler worker not configured")
		if err != nil {
			return err
		}
		return nil
	}

	reqURL, err := url.Parse(cfg.crawlerURL + "/stats")
	if err != nil {
		return err
	}
	query := reqURL.Query()
	parts := strings.Split(msg, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "keywords=") {
			keywordsQueryVal := strings.TrimPrefix(part, "keywords=")
			query.Set("keywords", keywordsQueryVal)
		}
		if strings.HasPrefix(part, "title=") {
			titleQueryVal := strings.TrimPrefix(part, "title=")
			query.Set("title", titleQueryVal)
		}
	}
	reqURL.RawQuery = query.Encode()

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return err
	}

	matches := []Match{}
	err = json.NewDecoder(resp.Body).Decode(&matches)
	if err != nil {
		return err
	}

	total := 0
	for _, match := range matches {
		total += match.Count
	}

	err = cfg.db.createMessage(crawlerBotAuthorUsername,
		fmt.Sprintf("Beep boop! I have found %d matches so far in books matching your request. Your raw query to me was: %s", total, reqURL.String()),
	)
	if err != nil {
		return err
	}

	return nil
}
