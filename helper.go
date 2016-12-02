package mux

// checkPairs returns the count of strings passed in, and an error if
import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// the count is not an even number.
func isEvenPairs(pairs ...string) (int, error) {
	length := len(pairs)
	if length%2 != 0 {
		return length, fmt.Errorf("mux: number of parameters must be multiple of 2, got %v", pairs)
	}
	return length, nil
}

// convertStringsToMap converts variadic string parameters to a
// string to string map.
func convertStringsToMapString(iep func(pairs ...string) (int, error), pairs ...string) (map[string]comparison, error) {

	buildComparator := func(pair string) (comparison, error) {
		return stringComparison(pair), nil
	}

	return genericConvertStringsToMap(iep, buildComparator, pairs...)
}

// mapFromPairsToRegex converts variadic string paramers to a
// string to regex map.
func convertStringsToMapRegex(iep func(pairs ...string) (int, error), pairs ...string) (map[string]comparison, error) {

	buildComparator := func(pair string) (comparison, error) {
		regex, err := regexp.Compile(pair)
		if err != nil {
			return nil, err
		}

		return regexComparsion{
			r: regex,
		}, nil
	}

	return genericConvertStringsToMap(iep, buildComparator, pairs...)
}

// genericConvertStringsToMap converts variadic string paramers to a
// string to whatever.
func genericConvertStringsToMap(
	iep func(pairs ...string) (int, error),
	buildComparator func(pair string) (comparison, error),
	pairs ...string) (map[string]comparison, error) {

	length, err := iep(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]comparison, length/2)
	for i := 0; i < length; i += 2 {

		cmp, err := buildComparator(pairs[i+1])

		if err != nil {
			return nil, err
		}

		m[pairs[i]] = cmp
	}
	return m, nil
}

type comparison interface {
	compare(string) bool
	isNotEmpty() bool
}

type stringComparison string

func (sc stringComparison) compare(value string) bool {
	if string(sc) == value {
		return true
	}
	return false
}

func (sc stringComparison) isNotEmpty() bool {
	return string(sc) != ""
}

type regexComparsion struct {
	r *regexp.Regexp
}

func (rc regexComparsion) compare(value string) bool {
	if rc.r.MatchString(value) {
		return true
	}
	return false
}

func (rc regexComparsion) isNotEmpty() bool {
	return rc.r != nil
}

// matchMapWithString returns true if the given key/value pairs exist in a given map.
func matchMap(compare map[string]comparison, toCompare map[string][]string, canonicalKey bool) bool {
	for k, v := range compare {
		// Check if key exists.
		if canonicalKey {
			k = http.CanonicalHeaderKey(k)
		}

		values := toCompare[k]

		if values == nil {
			return false
		}

		if v.isNotEmpty() {
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values {
				if v.compare(value) {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}

	return true
}

// containsRegexPath returns true if the path a regex path
func containsRegex(path string) bool {
	return strings.Contains(path, "#")
}

// containsRegexPath returns true if the path contains vars
func containsVars(path string) bool {
	return strings.Contains(path, ":")
}
