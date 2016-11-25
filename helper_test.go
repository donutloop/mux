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

func TestMapFromPairsToString(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 2, nil
	}

	pairs := []string{"content-type", "application/json"}

	m, _ := mapFromPairsToString(isEvenPairs, pairs...)

	if value, ok := m["content-type"]; !ok || value != "application/json" {
		t.Error("Unexpected pair")
	}
}

func TestMapFromPairsToStringError(t *testing.T) {
	isEvenPairs := func(pairs ...string) (int, error) {
		return 3, errors.New("Something went wrong")
	}

	pairs := []string{}

	if _, err := mapFromPairsToString(isEvenPairs, pairs...); err == nil {
		t.Error("Unexpected nil error")
	}
}
