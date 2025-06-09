package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

const maxResults = 10

func Top10(s string) []string {
	sSplit := strings.Fields(s)
	numStr := map[string]int{}
	for _, str := range sSplit {
		if str == "-" {
			continue
		}
		newStr := strings.TrimFunc(str, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		if newStr != "" {
			newStrLower := strings.ToLower(newStr)
			numStr[newStrLower]++
		} else {
			numStr[str]++
		}
	}

	words := make([]string, 0, len(numStr))
	for k := range numStr {
		words = append(words, k)
	}
	sort.Slice(words, func(i, j int) bool {
		a, b := words[i], words[j]
		if numStr[a] == numStr[b] {
			return a < b
		}
		return numStr[a] > numStr[b]
	})

	if len(words) > maxResults {
		return words[:maxResults]
	}
	return words
}
