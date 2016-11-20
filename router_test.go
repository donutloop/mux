package mux

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PathTest struct {
	title string
	path  string
}

func TestHost(t *testing.T) {

	tests := []PathTest{
		{
			title: "Path route with single path , match",
			path:  "/api/",
		},
	}

	for _, test := range tests {
		code, ok := testGET(test)

		if !ok {
			t.Errorf("Expected status code 200, Actucal status code %v", code)
		}
	}
}

func testGET(pt PathTest) (int, bool) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}
	r := NewRouter()

	r.Get(pt.path, handler)

	req, _ := http.NewRequest("GET", "http://localhost"+pt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, false
	}

	if res.Code != 200 && content.String() == "Hello World!" {
		return res.Code, false
	}

	return res.Code, true
}
