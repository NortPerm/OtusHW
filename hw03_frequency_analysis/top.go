package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const excludePunctSigns = `-`

var (
	removePunct    = regexp.MustCompile(`[^[:^punct:]` + excludePunctSigns + `]`)
	removeAllPunct = regexp.MustCompile(`[[:punct:]]`)
)

type analyzer struct {
	words       []string
	uniqueWords []string
	freqMap     map[string]int
}

func newAnalyzer(inputText string) *analyzer {
	return &analyzer{removePuncSigns(strings.Fields(inputText)), []string{}, make(map[string]int)}
}

func (a *analyzer) buildFreqMap() {
	for _, word := range a.words {
		if word != "" {
			if a.freqMap[word] == 0 {
				// if this word has not been encountered before, add it to the list of unique
				a.uniqueWords = append(a.uniqueWords, word)
			}
			a.freqMap[word]++
		}
	}
}

func removePuncSigns(slice []string) []string {
	// for non-asterix task just comment the code except last line (return slice)
	for ind, word := range slice {
		// if after removing all punctuatioan we have empty string
		if removeAllPunct.ReplaceAllString(strings.ToLower(word), "") != "" {
			// we remove all punctuation except dash
			slice[ind] = removePunct.ReplaceAllString(strings.ToLower(word), "")
		} else {
			// otherwise set word as empty string
			slice[ind] = ""
		}
	}
	return slice
}

func cut10(slice []string) []string {
	if len(slice) < 10 {
		return slice
	}
	return slice[:10]
}

func Top10(inputText string) []string {
	analyzer := newAnalyzer(inputText)
	analyzer.buildFreqMap()
	sort.Slice(analyzer.uniqueWords, func(i, j int) bool {
		freqI := analyzer.freqMap[analyzer.uniqueWords[i]]
		freqJ := analyzer.freqMap[analyzer.uniqueWords[j]]
		if freqI == freqJ {
			return analyzer.uniqueWords[i] < analyzer.uniqueWords[j]
		}
		return freqI > freqJ
	})
	return cut10(analyzer.uniqueWords)
}
