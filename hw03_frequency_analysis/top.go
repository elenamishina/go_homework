package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var punctRegex = regexp.MustCompile(`^\p{P}+|\p{P}+$`)

func Top10(s string) []string {
	sSplit := strings.Fields(s)
	numStr := map[string]int{}
	for _, str := range sSplit {
		if str == "-" {
			continue
		}
		newStr := punctRegex.ReplaceAllString(str, "")
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
		if numStr[words[i]] != numStr[words[j]] {
			return numStr[words[i]] > numStr[words[j]]
		}
		return words[i] < words[j]
	})

	num := 10
	if len(words) < num {
		num = len(words)
	}
	return words[:num]
}
