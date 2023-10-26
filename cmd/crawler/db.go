package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"unicode"
)

type DB interface {
	getCounts() ([]Match, error)
	saveCount(title string, count int) error
}

type Disk struct {
	crawlerDBPath string
}

func (d Disk) getCounts() ([]Match, error) {
	entries, err := os.ReadDir(d.crawlerDBPath)
	if err != nil {
		return nil, err
	}

	matchesSlice := []Match{}
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
		matchesSlice = append(matchesSlice, match)
	}
	slices.SortFunc(matchesSlice, func(a, b Match) int {
		return a.Count - b.Count
	})
	return matchesSlice, err
}

func (d Disk) saveCount(title string, count int) error {
	cleanedTitle := title
	cleanedTitle = strings.ToLower(title)
	cleanedTitle = removeSpaces(cleanedTitle)
	cleanedTitle = removePunctuation(cleanedTitle)
	match := Match{
		BookTitle: title,
		Count:     count,
	}
	dat, err := json.Marshal(match)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d.crawlerDBPath, cleanedTitle+".json"), dat, 0755)
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
	matches map[string]Match
}

func (m *Memory) getCounts() ([]Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	matchesSlice := []Match{}
	for _, match := range m.matches {
		matchesSlice = append(matchesSlice, match)
	}
	slices.SortFunc(matchesSlice, func(a, b Match) int {
		return a.Count - b.Count
	})
	return matchesSlice, nil
}

func (m *Memory) saveCount(title string, count int) error {
	cleanedTitle := title
	cleanedTitle = strings.ToLower(title)
	cleanedTitle = removeSpaces(cleanedTitle)
	cleanedTitle = removePunctuation(cleanedTitle)
	match := Match{
		BookTitle: title,
		Count:     count,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.matches[cleanedTitle] = match
	return nil
}
