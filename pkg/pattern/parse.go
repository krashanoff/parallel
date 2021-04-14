package pattern

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	indexFmt = `([^:]+)?:(-?[0-9]+)(:(-?[0-9]+))?`
)

var (
	indexCheck *regexp.Regexp
)

func init() {
	indexCheck = regexp.MustCompile(indexFmt)
}

func parseBrackets(bracket, replacement string) (string, error) {
	if len(bracket) < 2 {
		return "", errors.New("bracket patterns must be wrapped in '{}'")
	}

	if bracket[0] != '{' || bracket[len(bracket)-1] != '}' {
		return "", errors.New("input pattern does not start with '{' and end with '}'")
	}

	// Literal insertion of '{}'.
	if len(bracket) == 2 {
		return replacement, nil
	}

	// Pattern-based insertion.
	pattern := bracket[1 : len(bracket)-1]
	if pattern[0] == '{' && pattern[len(pattern)-1] == '}' {
		return pattern, nil
	}

	// Part-based index replacement.
	matches := indexCheck.FindStringSubmatch(pattern)
	splitBy := matches[1]
	replaceSplit := strings.Split(replacement, splitBy)

	// Parse start index.
	startIdx, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		startIdx = 0
	}

	// Negative index handling.
	if startIdx < int64(0) {
		startIdx += int64(len(replaceSplit))
	}

	endIdx, err := strconv.ParseInt(matches[4], 10, 64)
	if err != nil {
		endIdx = startIdx + 1
	}

	log.Println(startIdx, endIdx)

	return strings.Join(replaceSplit[startIdx:endIdx], ""), nil
}

/* InsertFilename replaces instances of {}-type patterns with
 * the provided syntax:
 *
 * - {} inserts the entire string.
 * - {{S}} inserts the string "{S}".
 * - {:n} references the nth character of the string.
 * - {:n:m} references the n through mth characters of the
 *   string.
 * - {S:n} references the nth component of the filename using
 *   'S' as a delimiter pattern. Colons are disallowed.
 * - {S:n:m} references the n through mth components of the
 *   filename using 'S' as a delimiter pattern. Colons are
 *   disallowed.
 *
 * Ill-formatted input pattern strings return an error.
 */
func InsertFilename(pattern, replacement string) (result string, err error) {
	stack := 0
	lastStart := -1

	for idx, c := range pattern {
		switch c {
		case '{':
			if stack == 0 {
				lastStart = idx
			}
			stack++
		case '}':
			stack--
			if stack == 0 {
				if s, err := parseBrackets(pattern[lastStart:idx+1], replacement); err != nil {
					return "", err
				} else {
					result += s
				}
			}
		default:
			if stack == 0 {
				result += string(c)
			}
		}
	}

	return
}

// GeneratePatterns generates a list of patterned strings using
// the input set replaceWith.
//
// For a replaceWith list of n strings, n strings are generated
// as output.
func GeneratePatterns(pattern string, replacements []string) ([]string, error) {
	result := make([]string, 0)

	for _, replacement := range replacements {
		if inserted, err := InsertFilename(pattern, replacement); err != nil {
			return nil, err
		} else {
			result = append(result, inserted)
		}
	}

	return result, nil
}
