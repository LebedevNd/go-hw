package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func Top10(text string) []string {
	var words []string
	var wordsMap map[string]int
	wordsResult := make([]string, 0, 10)
	wordsMap = make(map[string]int)

	text = prepareText(text)
	words = strings.Fields(text)

	wordsCounts := make([]wordCount, 0, len(words))

	for _, word := range words {
		word = strings.ToLower(word)
		if word != "" {
			wordsMap[word]++
		}
	}

	for word, count := range wordsMap {
		wordsCounts = append(wordsCounts, wordCount{word, count})
	}

	sort.Slice(wordsCounts, func(i, j int) bool {
		if wordsCounts[i].count != wordsCounts[j].count {
			return wordsCounts[i].count > wordsCounts[j].count
		}
		return wordsCounts[i].word < wordsCounts[j].word
	})

	for _, wordCount := range wordsCounts {
		wordsResult = append(wordsResult, wordCount.word)
		if len(wordsResult) == 10 {
			break
		}
	}

	return wordsResult
}

func prepareText(text string) string {
	re := regexp.MustCompile(`[[:punct:]]`)
	text = re.ReplaceAllString(text, "")
	return text
}
