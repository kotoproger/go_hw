package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	re       = regexp.MustCompile(`^[[:punct:]]+|[[:punct:]]+$`)
	dashWord = regexp.MustCompile(`^[-]{2,}$`)
)

func Top10(input string) []string {
	words := make(map[string]int)
	wordsSlice := []string{}

	for _, word := range strings.Fields(input) {
		lowerCaseWord := strings.ToLower(word)
		var lowerCaseWithoutPuncts string

		if !dashWord.Match([]byte(lowerCaseWord)) {
			lowerCaseWithoutPuncts = string(re.ReplaceAll([]byte(lowerCaseWord), []byte{}))
		} else {
			lowerCaseWithoutPuncts = lowerCaseWord
		}

		if lowerCaseWithoutPuncts == "" {
			continue
		}

		words[lowerCaseWithoutPuncts]++
		if words[lowerCaseWithoutPuncts] == 1 {
			wordsSlice = append(wordsSlice, lowerCaseWithoutPuncts)
		}
	}
	if len(wordsSlice) == 0 {
		return wordsSlice
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		return words[wordsSlice[i]] > words[wordsSlice[j]] ||
			(words[wordsSlice[i]] == words[wordsSlice[j]] && wordsSlice[i] < wordsSlice[j])
	})

	if len(wordsSlice) < 10 {
		return wordsSlice
	}
	return wordsSlice[:10]
}
