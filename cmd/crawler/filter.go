package main

import (
	"strings"
)

func filterKeywords(filterKeywordsString string, matches map[string]map[string]Match) map[string]map[string]Match {
	filterKeywords := strings.Split(filterKeywordsString, ",")
	filterKeywordsSet := map[string]struct{}{}
	for _, keyword := range filterKeywords {
		filterKeywordsSet[keyword] = struct{}{}
	}
	if len(filterKeywordsSet) == 0 {
		return matches
	}
	filterMatches := map[string]map[string]Match{}
	for keyword, matches := range matches {
		if _, ok := filterKeywordsSet[keyword]; ok {
			filterMatches[keyword] = matches
		}
	}
	return filterMatches
}

func filterTitles(filterTitleString string, matches map[string]map[string]Match) map[string]map[string]Match {
	if len(filterTitleString) == 0 {
		return matches
	}
	filterMatches := map[string]map[string]Match{}
	for keyword, kMatches := range matches {
		for title, match := range kMatches {
			if strings.Contains(strings.ToLower(title), strings.ToLower(filterTitleString)) {
				if _, ok := filterMatches[keyword]; !ok {
					filterMatches[keyword] = map[string]Match{}
				}
				filterMatches[keyword][title] = match
			}
		}
	}
	return filterMatches
}
