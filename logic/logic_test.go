package logic

import (
	"fmt"
	"log"
	td "sanitize/testdata"
	"slices"
	"sort"
	"strings"
	"testing"
)

var unitTest = []string{
	"kBD6 <string0> qGWpWk", "k <string0> BD6 <string1> qGWpWk", "k <string0> BD6 <string2> qGWp <string3> Wk", "kB D6 qGWpWk <string0>"}

func TestFilter(t *testing.T) {
	var requests []string
	var expected []string

	for _, str := range td.Data {
		requests = append(requests, strings.Replace(unitTest[0], "<string0>", str, -1))
		expected = append(expected, strings.Replace(unitTest[0], "<string0>", strings.Repeat("*", len(str)), -1))
	}

	result, err := SanitizeText(requests, td.Data)
	if err != nil {
		t.Error(err)
	}

	if !slices.Equal(result, expected) {
		for index, res := range result {
			if res != expected[index] {
				log.Print(fmt.Sprintf("%s != %s\n", res, expected[index]))
			}
		}

		t.Error("Not complying to expected output")
	}
}

func TestFilterTwo(t *testing.T) {
	var requests []string
	var expected []string

	for _, str := range td.Data {
		req := unitTest[0]
		exp := unitTest[0]
		for i := 0; i < 2; i++ {
			req = strings.Replace(req, fmt.Sprintf("<string%d>", i), str, -1)
			exp = strings.Replace(exp, fmt.Sprintf("<string%d>", i), strings.Repeat("*", len(str)), -1)
		}
		requests = append(requests, req)
		expected = append(expected, exp)
	}

	result, err := SanitizeText(requests, td.Data)
	if err != nil {
		t.Error(err)
	}

	if !slices.Equal(result, expected) {
		for index, res := range result {
			if res != expected[index] {
				log.Print(fmt.Sprintf("%s != %s\n", res, expected[index]))
			}
		}

		t.Error("Not complying to expected output")
	}
}

func TestFilterThree(t *testing.T) {
	var requests []string
	var expected []string

	for _, str := range td.Data {
		req := unitTest[0]
		exp := unitTest[0]
		for i := 0; i < 3; i++ {
			req = strings.Replace(req, fmt.Sprintf("<string%d>", i), str, -1)
			exp = strings.Replace(exp, fmt.Sprintf("<string%d>", i), strings.Repeat("*", len(str)), -1)
		}
		requests = append(requests, req)
		expected = append(expected, exp)
	}

	result, err := SanitizeText(requests, td.Data)
	if err != nil {
		t.Error(err)
	}

	if !slices.Equal(result, expected) {
		for index, res := range result {
			if res != expected[index] {
				log.Print(fmt.Sprintf("%s != %s\n", res, expected[index]))
			}
		}

		t.Error("Not complying to expected output")
	}
}

func TestFilterWhitespace(t *testing.T) {
	var requests []string
	var expected []string

	for _, str := range td.Data {
		requests = append(requests, strings.Replace(unitTest[3], "<string0>", str, -1))
		expected = append(expected, strings.Replace(unitTest[3], "<string0>", strings.Repeat("*", len(str)), -1))
	}

	result, err := SanitizeText(requests, td.Data)
	if err != nil {
		t.Error(err)
	}

	if !slices.Equal(result, expected) {
		for index, res := range result {
			if res != expected[index] {
				log.Print(fmt.Sprintf("%s != %s\n", res, expected[index]))
			}
		}

		t.Error("Not complying to expected output")
	}
}

func TestFilterTwoWords(t *testing.T) {
	sanitizeWords := td.Data
	sort.Slice(sanitizeWords, func(i, j int) bool {
		return len(sanitizeWords[i]) > len(sanitizeWords[j])
	})

	var requests []string
	var expected []string

	for i := 1; i <= len(sanitizeWords)/2; i++ {
		firstWord := sanitizeWords[i-1]
		secondWord := sanitizeWords[((len(sanitizeWords)/2)+i)-1]

		requests = append(requests, strings.Replace(strings.Replace(unitTest[1], "<string0>",
			firstWord, -1), "<string1>", secondWord, -1))

		expected = append(expected, strings.Replace(
			strings.Replace(unitTest[1], "<string0>", strings.Repeat("*", len(firstWord)), -1),
			"<string1>", strings.Repeat("*", len(secondWord)), -1))
	}

	result, err := SanitizeText(requests, td.Data)
	if err != nil {
		t.Error(err)
	}

	if !slices.Equal(result, expected) {
		for index, res := range result {
			if res != expected[index] {
				log.Print(fmt.Sprintf("%s != %s\n", res, expected[index]))
			}
		}

		t.Error("Not complying to expected output")
	}
}
