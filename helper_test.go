package mux

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

func createRandomStringsCount() (pairs []string, count int) {
	randNumber := rand.Intn(100-10) + 10
	for i := 0; i < randNumber; i++ {
		pairs = append(pairs, "dummy")
	}
	count = len(pairs)

	return
}

func TestIsEvenPairs(t *testing.T) {
	for i := 0; i <= 5; i++ {
		pairs, count := createRandomStringsCount()
		t.Run(fmt.Sprintf("Check pairs test (Count: %v)", count), func(t *testing.T) {
			_, err := isEvenPairs(pairs...)

			if err == nil && count%2 != 0 {
				t.Error("Unexpected pairs count")
			} else if err != nil && count%2 == 0 {
				t.Error("Unexpected pairs count")
			}
		})
	}
}

func TestConvertStringsToMapString(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 2, nil
	}

	pairs := []string{"content-type", "application/json"}

	m, _ := convertStringsToMapString(isEvenPairs, pairs...)

	if value, ok := m["content-type"]; !ok || !value.compare("application/json") {
		t.Errorf("Unexpected pair (%s)", string(value.(stringComparison)))
	}
}

func TestConvertStringsToMapRegex(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 2, nil
	}

	pairs := []string{"content-type", "application/json"}

	m, _ := convertStringsToMapRegex(isEvenPairs, pairs...)

	if value, ok := m["content-type"]; !ok || !value.compare("application/json") {
		t.Errorf("Unexpected pair (%s)", value.(regexComparsion))
	}
}

func TestGenericConvertStringsToMapFail(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 2, nil
	}

	buildComparator := func(pair string) (comparison, error) {
		return nil, errors.New("Somthing went wrong")
	}

	pairs := []string{"content-type", "application/json"}

	if _, err := genericConvertStringsToMap(isEvenPairs, buildComparator, pairs...); err == nil {
		t.Error("Expected a error")
	}
}

func TestConvertStringsToMapError(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 3, errors.New("Something went wrong")
	}

	pairs := []string{}

	if _, err := convertStringsToMapString(isEvenPairs, pairs...); err == nil {
		t.Error("Unexpected nil error")
	}
}

func TestMatchMap(t *testing.T) {

	compare := map[string]comparison{
		"content-type": stringComparison("application/json"),
	}

	toCompare := map[string][]string{
		"content-type": {
			"application/json",
		},
	}

	ok := matchMap(compare, toCompare, false)

	if !ok {
		t.Error("Unexpected non match")
	}
}

func BenchmarkMatchMap(b *testing.B) {

	tests := []struct {
		title   string
		compare map[string]comparison
	}{
		{
			title: "MatchMap (String comparison): 1 Value",
			compare: map[string]comparison{
				"content-type": stringComparison("application/json"),
			},
		},
		{
			title: "MatchMap (String comparison): 2 Values",
			compare: map[string]comparison{
				"accept-encoding": stringComparison("gzip"),
				"content-type":    stringComparison("application/json"),
			},
		},
		{
			title: "MatchMap (String comparison): 3 Values",
			compare: map[string]comparison{
				"accept-encoding": stringComparison("gzip"),
				"content-type":    stringComparison("application/json"),
				"accept":          stringComparison("text/plain"),
			},
		},
	}

	toCompare := map[string][]string{
		"content-type": {
			"applcation/json",
		},
		"accept-charset": {
			"utf-8",
		},
		"accept-encoding": {
			"gzip",
			"deflate",
		},
		"accept-language": {
			"en-US",
		},
		"cache-control": {
			"no-cache",
		},
		"date": {
			"Date: Tue, 15 Nov 1994 08:12:31 GMT",
		},
		"max-Forwards": {
			"10",
		},
		"accept": {
			"text/plain",
			"text/html",
		},
	}

	for _, test := range tests {
		b.Run(test.title, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				matchMap(test.compare, toCompare, false)
			}
		})
	}
}
