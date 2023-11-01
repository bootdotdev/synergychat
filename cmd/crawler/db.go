package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type Match struct {
	Keyword   string
	BookTitle string
	Count     int
}

type DB interface {
	// keyword -> title -> count
	getCounts() (map[string]map[string]Match, error)
	saveCount(title, keyword string, count int) error
	init() error
}

type Disk struct {
	crawlerDBPath string
}

func (d *Disk) getCounts() (map[string]map[string]Match, error) {
	entries, err := os.ReadDir(d.crawlerDBPath)
	if err != nil {
		return nil, err
	}
	matches := map[string]map[string]Match{}
	for _, entry := range entries {
		dat, err := os.ReadFile(filepath.Join(d.crawlerDBPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		match := Match{}
		err = json.Unmarshal(dat, &match)
		if err != nil {
			return nil, err
		}
		if _, ok := matches[match.Keyword]; !ok {
			matches[match.Keyword] = map[string]Match{}
		}
		matches[match.Keyword][match.BookTitle] = match
	}
	return matches, nil
}

func (d *Disk) saveCount(keyword, title string, count int) error {
	cleanedTitle := title
	cleanedTitle = strings.ToLower(title)
	cleanedTitle = removeSpaces(cleanedTitle)
	cleanedTitle = removePunctuation(cleanedTitle)
	match := Match{
		BookTitle: title,
		Count:     count,
		Keyword:   keyword,
	}
	dat, err := json.Marshal(match)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d.crawlerDBPath, keyword+"."+cleanedTitle+".json"), dat, 0755)
}

func (d *Disk) init() error {
	err := os.MkdirAll(d.crawlerDBPath, 0755)
	return err
}

func removePunctuation(s string) string {
	var result strings.Builder
	for _, r := range s {
		if !unicode.IsPunct(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func removeSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

type Memory struct {
	mu      *sync.Mutex
	matches map[string]map[string]Match
}

func (m *Memory) getCounts() (map[string]map[string]Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return deepCopy(m.matches), nil
}

func (m *Memory) saveCount(keyword, title string, count int) error {
	cleanedTitle := title
	cleanedTitle = strings.ToLower(title)
	cleanedTitle = removeSpaces(cleanedTitle)
	cleanedTitle = removePunctuation(cleanedTitle)
	match := Match{
		BookTitle: title,
		Count:     count,
		Keyword:   keyword,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.matches[keyword]
	if !ok {
		m.matches[keyword] = map[string]Match{}
	}
	m.matches[keyword][cleanedTitle] = match
	return nil
}

func (m *Memory) init() error {
	if m != nil {
		return nil
	}
	m = &Memory{
		mu:      &sync.Mutex{},
		matches: map[string]map[string]Match{},
	}
	return nil
}

func deepCopy(m map[string]map[string]Match) map[string]map[string]Match {
	result := map[string]map[string]Match{}
	for k, v := range m {
		result[k] = map[string]Match{}
		for a, b := range v {
			result[k][a] = b
		}
	}
	return result
}

func matchesMapToSlice(m map[string]map[string]Match) []Match {
	result := []Match{}
	for _, v := range m {
		for _, b := range v {
			result = append(result, b)
		}
	}
	return result
}
