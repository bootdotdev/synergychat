package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilterKeywords(t *testing.T) {
	tests := []struct {
		filterKeywordsStr string
		matches           map[string]map[string]Match
		want              map[string]map[string]Match
	}{
		{
			filterKeywordsStr: "technology",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
			},
		},
		{
			filterKeywordsStr: "technology,science",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
		},

		{
			filterKeywordsStr: "fiction",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			if got := filterKeywords(tt.filterKeywordsStr, tt.matches); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterKeywords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterTitles(t *testing.T) {
	tests := []struct {
		filterTitleString string
		matches           map[string]map[string]Match
		want              map[string]map[string]Match
	}{
		{
			filterTitleString: "Article 1",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
			},
		},
		{
			filterTitleString: "article 1",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
			},
		},
		{
			filterTitleString: "article",
			matches: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
			want: map[string]map[string]Match{
				"technology": {"Article 1": {Keyword: "technology", BookTitle: "Article 1", Count: 10}},
				"health":     {"Article 2": {Keyword: "health", BookTitle: "Article 2", Count: 5}},
				"science":    {"Article 3": {Keyword: "science", BookTitle: "Article 3", Count: 7}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			if got := filterTitles(tt.filterTitleString, tt.matches); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterTitles() = %v, want %v", got, tt.want)
			}
		})
	}
}
