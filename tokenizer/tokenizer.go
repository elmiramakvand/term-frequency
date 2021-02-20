package tokenizer

import (
	"regexp"
	"strings"
)

type queryStrings []string

// standard tokenizer : This tokenizer splits the text field into tokens,
// treating whitespace and punctuation as delimiters. Delimiter characters are discarded
func StandardTokenizer(qs queryStrings) []string {
	var tokens []string
	for _, queryString := range qs {
		step1 := strings.ToLower(queryString)
		var re = regexp.MustCompile(`(^\.*)| \.| *\. |@|'|\?|\(|\)|"|“|”|,|-|:|\.*$`)
		step2 := re.ReplaceAllString(step1, " ")
		tokens = append(tokens, strings.Fields(step2)...)
	}
	return tokens
}

// keyword tokenizer : This tokenizer treats the entire text field as a single token.
func KeywordTokenizer(qs queryStrings, tokens []string) []string {
	var result []string
	for _, queryString := range qs {
		found := Find(tokens, strings.ToLower(queryString))
		if !found {
			//Value not found in slice
			result = append(result, strings.ToLower(queryString))
		}
	}
	return result
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
