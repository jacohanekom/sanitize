package logic

import (
	"errors"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
)

// SanitizeText sanitizes the input text according to the provided word list and returns the result in the same order
// it was provided. In the case of an error the method will return a empty string list, and the appropriate error
func SanitizeText(textToSanitize []string, sanitizeWords []string) (result []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = errors.New("an error occurred, while sanitizing text")
		}
	}()

	//The slice had to be sorted due to the longest words first
	//In the case that it is not it will produce invalid results and cause uncensored text
	sort.Slice(sanitizeWords, func(i, j int) bool {
		return len(sanitizeWords[i]) > len(sanitizeWords[j])
	})

	for _, text := range textToSanitize {
		//Building a response string, and censoring as we continue to loop through the words
		//to sensor the requested string

		var replace []int
		for _, word := range sanitizeWords {
			//escaping special regex strings
			word = strings.ReplaceAll(word, "*", "\\*")

			re := regexp.MustCompile("\\b" + strings.ToUpper(word) + "\\b")
			matches := re.FindAllStringIndex(strings.ToUpper(text), -1)

			for _, match := range matches {
				for idx := match[0]; idx < match[1]; idx++ {
					replace = append(replace, idx)
				}
			}
		}

		sanitized := []rune(text)
		for _, index := range replace {
			if index < len(sanitized) {
				sanitized[index] = '*'
			}
		}

		result = append(result, string(sanitized))

	}

	return result, nil
}
