package tokenizer

import (
	"regexp"
	"strings"
)

// standard tokenizer : This tokenizer splits the text field into tokens,
// treating whitespace and punctuation as delimiters. Delimiter characters are discarded
func StandardTokenizer(qs string) []string {
	var tokens []string
	step1 := strings.ToLower(qs)
	var re = regexp.MustCompile(`(^\.*)| \.| *\. |@|'|\?|\(|\)|"|“|”|,|-|:|\.*$`)
	step2 := re.ReplaceAllString(step1, " ")
	tokens = append(tokens, strings.Fields(step2)...)
	return tokens
}

// keyword tokenizer : This tokenizer treats the entire text field as a single token.
func KeywordTokenizer(qs string, tokens []string) []string {
	var result []string

	found := Find(tokens, strings.ToLower(qs))
	if !found {
		//Value not found in slice
		result = append(result, strings.TrimSpace(strings.ToLower(qs)))
	}
	return result
}

func CheckStringHasWord(qs string) bool {
	var re = regexp.MustCompile(`(^\.*)| \.| *\. |@|'|\?|\(|\)|"|“|”|,|-|:|\.*$`)
	trimedQueryString := re.ReplaceAllString(strings.ToLower(qs), " ")
	if len(strings.ReplaceAll(trimedQueryString, " ", "")) == 0 {
		return false
	}
	return true
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
