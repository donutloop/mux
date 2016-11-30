[![Build Status](https://travis-ci.org/donutloop/mux.svg?branch=master)](https://travis-ci.org/donutloop/mux)

# What is mux ?

Mux is a lightweight and fast HTTP Multiplexer for Golang >= 1.7

Status: Alpha (Not ready for production)

## Features:

* REGEX URL Matcher
* Vars URl Matcher
* URL Matcher
* Header Matcher
* Scheme Matcher 
* Custom Matcher
* Route Validators 
* Http method declaration
* Support for standard lib http.Handler and http.HandlerFunc
* Custom NotFound handler
* Respect the Go standard http.Handler interface
* Routes are sorted
* Context support

## Example (Method GET):

```go
    package main

    import (
        "net/http"
        "fmt"
        "os"

        "github.com/donutloop/mux"
    )

    func main() {
        r := newRouter()

        r.HandleFunc(http.MethodGet, "/home", homeHandler)
        r.Handler(http.MethodGet, "/home-1", http.HandlerFunc(homeHandler)
        r.Get("/home-2", homeHandler)

        errs := r.ListenAndServe(":8080")

        for _ , err := range errs {
            fmt.print(err)
        }

        if 0 != len(errs) {
            os.Exit(2)
        }
    }

    func homeHandler(rw http.ResponseWriter, req *http.Request){
        //...
        rw.Write([]byte("Hello World!")
    }
```

## More documentaion comming soon