package mux

import (
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

func TestCheckPairs(t *testing.T) {

	for i := 0; i <= 5; i++ {
		pairs, count := createRandomStringsCount()
		t.Run(fmt.Sprintf("Check pairs test (Count: %v)", count), func(t *testing.T) {
			_, err := checkPairs(pairs...)

			if err == nil && count%2 != 0 {
				t.Error("Unexpected pairs count")
			} else if err != nil && count%2 == 0 {
				t.Error("Unexpected pairs count")
			}
		})
	}
}
