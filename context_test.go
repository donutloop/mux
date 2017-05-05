package mux

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetQueriesFail(t *testing.T) {
	r := new(http.Request)
	if value := GetQueries(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}

func TestCurrentRouteFail(t *testing.T) {
	r := new(http.Request)
	if value := CurrentRoute(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}

func TestGetVarsFail(t *testing.T) {
	r := new(http.Request)
	if value := GetVars(r); value != nil {
		t.Errorf("Unexpected value (%v)", value)
	}
}

func BenchmarkExtractQueries(b *testing.B) {
	request := &http.Request{
		URL: &url.URL{
			RawQuery: "limit=10&offset=10&gender=female&age[0]=20&age[0]=50",
		},
	}

	for n := 0; n < b.N; n++ {
		extractQueries(request)
	}
}
