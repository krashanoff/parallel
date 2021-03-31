package pkg

import (
	"errors"
	"regexp"
	// "strings"
)

const patternFmt = `.*\{((\{.*\})|(.?:\-?[0-9]+:(\-?[0-9]+)?))\}`
var patternCheck *regexp.Regexp

func init() {
	patternCheck = regexp.MustCompile(patternFmt)
}

// ValidPattern checks whether the given
// pattern adheres to parallel's syntax.
func ValidPattern(pattern string) bool {
	return patternCheck.Match([]byte(pattern))
}

/* InsertFilename replaces instances of {}-type patterns with
 * the provided syntax:
 *
 * - {} inserts the entire string.
 * - {{S}} inserts the string "{S}".
 * - {:n} references the nth character of the string.
 * - {:n:m} references the n through mth characters of the
 *   string.
 * - {c:n} references the nth component of the filename using
 *   'c' as a delimiter pattern.
 * - {c:n:m} references the n through mth components of the
 *   filename using 'c' as a delimiter pattern.
 * - {n:m} references components [n, m).
 *
 * Ill-formatted input pattern strings return an error.
 */
func InsertFilename(pattern string, replacement string) (string, error) {
	if !ValidPattern(pattern) {
		return "", errors.New("input pattern is ill-formatted")
	}

	for _, c := range pattern {
		switch c {
		case '{':
		case '}':
		case ':':
		}
	}

	return "", nil
}

// GeneratePatterns generates a list of patterned strings using
// the input set replaceWith.
//
// For a replaceWith list of n strings, n strings are generated
// as output.
func GeneratePatterns(pattern string, replacements []string) ([]string, error) {
	return nil, nil
}