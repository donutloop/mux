package mux

import (
	"net/http"
	"sort"
	"strings"
	"testing"
)

func TestNewBadRouteError(t *testing.T) {
	r := &Route{
		methodName: "GET",
		path:       "/api/user",
	}

	err := newBadRouteError(r, "Something went wrong")

	if err == nil || !strings.Contains(err.Error(), "Something went wrong") {
		t.Errorf("Bad error message (%v)", err.Error())
	}
}

func TestSortMatchers(t *testing.T) {
	matcherFunc := func(*http.Request) bool {
		return true
	}
	mf := MatcherFunc(matcherFunc)

	ms := matchers([]Matcher{})
	ms = append(ms, newSchemeMatcher("https"), pathMatcher("/api/"), mf)
	sort.Sort(ms)

	if ms[0].Rank() != rankAny || ms[1].Rank() != rankPath || ms[2].Rank() != rankScheme {
		t.Errorf("Unexpected ranking (Index 0: %d, Index 1: %d, Index 2: %d)", ms[0].Rank(), ms[1].Rank(), ms[2].Rank())
	}
}
