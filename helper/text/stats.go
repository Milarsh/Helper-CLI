package text

import (
	_ "embed"
	"regexp"
	"slices"
	"strings"
)

//go:embed stop-words.txt
var stopWordsFile string

var nonWord = regexp.MustCompile(`[^\p{L}\p{N}\s]+`)
var stopWords = strings.Fields(strings.TrimSpace(stopWordsFile))

func CountSymbols(str string) int {
	return len(str)
}

func Tokenize(str string) []string {
	text := strings.ToLower(str)

	text = nonWord.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	return strings.Fields(text)
}

func CountWords(str string) int {
	return len(Tokenize(str))
}

func StopWordsList() []string {
	return stopWords
}

func CountStopWords(str string) int {
	counter := 0
	words := Tokenize(str)
	for _, word := range words {
		if slices.Contains(StopWordsList(), word) {
			counter++
		}
	}
	return counter
}
