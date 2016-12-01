[![Build Status](https://travis-ci.org/donutloop/mux.svg?branch=master)](https://travis-ci.org/donutloop/mux)

# What is mux ?

Mux is a lightweight and fast HTTP Multiplexer for Golang >= 1.7

Status: Alpha (Not ready for production)

## Features:

* REGEX URL Matcher
* Vars URL Matcher
* GetVars in handler
* GetQueries in handler
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

## Roadmap

* HeaderRegex Matcher

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
        r := mux.NewRouter()

        r.HandleFunc(http.MethodGet, "/home", homeHandler)
        
        r.Handler(http.MethodGet, "/home-1", http.HandlerFunc(homeHandler))
        
        r.Get("/home-2", homeHandler)
        
        r.RegisterRoute(http.MethodGet, r.NewRoute().Path("/home-3").HandlerFunc(homeHandler))
        
        r.RegisterRoute(http.MethodGet, r.NewRoute().Path("/home-4").Handler(http.HandlerFunc(homeHandler)))
        
    	errorHandler := func(errs []error) {
            for _ , err := range errs {
                fmt.Print(err)
            }
            if 0 != len(errs) {
                os.Exit(2)
            }
	    }

        errs := r.ListenAndServe(":8080", errorHandler)
    }

    func homeHandler(rw http.ResponseWriter, req *http.Request) {
        //...
        rw.Write([]byte("Hello World!")
    }
```

## Example (Method POST):

```go
    package main

    import (
        "net/http"
        "fmt"
        "os"

        "github.com/donutloop/mux"
    )

    func main() {
        r := mux.NewRouter()

        r.HandleFunc(http.MethodPost, "/user/create", userHandler)
        
        r.Handler(http.MethodPost, "/user-1/create", http.HandlerFunc(userHandler)
        
        r.Post("/user-2/create", userHandler)

        r.RegisterRoute(http.MethodPost, r.NewRoute().Path("/user-3/create").HandlerFunc(userHandler))
        
        r.RegisterRoute(http.MethodPost, r.NewRoute().Path("/user-4/create").Handler(http.HandlerFunc(userHandler)))
        
    	errorHandler := func(errs []error) {
            for _ , err := range errs {
                fmt.Print(err)
            }
            if 0 != len(errs) {
                os.Exit(2)
            }
	    }

        errs := r.ListenAndServe(":8080", errorHandler)
    }

    func userHandler(rw http.ResponseWriter, req *http.Request) {
        //...
        rw.Write([]byte("Created successfully a new user")
    }
```

## Example (Method GET & Scheme Matcher):

```go
    package main

    import (
        "net/http"
        "fmt"
        "os"

        "github.com/donutloop/mux"
    )

    func main() {
        r := mux.NewRouter()
        
        r.Get("/home", homeHandler).Schemes("https")
        
    	errorHandler := func(errs []error) {
            for _ , err := range errs {
                fmt.Print(err)
            }
            if 0 != len(errs) {
                os.Exit(2)
            }
	    }

        errs := r.ListenAndServe(":8080", errorHandler)
    }

    func homeHandler(rw http.ResponseWriter, req *http.Request) {
        //...
        rw.Write([]byte("Hello world")
    }
```
## Example (Method Put & GetVars):

```go
    package main

    import (
        "net/http"
        "fmt"
        "os"

        "github.com/donutloop/mux"
    )

    func main() {
        r := mux.NewRouter()
        
        r.Post("/user/update/:number", userHandler)

    	errorHandler := func(errs []error) {
            for _ , err := range errs {
                fmt.Print(err)
            }
            if 0 != len(errs) {
                os.Exit(2)
            }
	    }

        errs := r.ListenAndServe(":8080", errorHandler)
    }

    func userHandler(rw http.ResponseWriter, req *http.Request) {
        userId := GetVars(req).Get(":number")
        //...
        rw.Write([]byte("Updated successfully a new user")
    }
```

## Example (Method GET & GetQueries):

```go
    package main

    import (
        "net/http"
        "fmt"
        "os"

        "github.com/donutloop/mux"
    )

    func main() {
        r := mux.NewRouter()
        
        //URL: http://localhost:8080/users?limit=10
        r.Post("/users", userHandler)

    	errorHandler := func(errs []error) {
            for _ , err := range errs {
                fmt.Print(err)
            }
            if 0 != len(errs) {
                os.Exit(2)
            }
	    }

        errs := r.ListenAndServe(":8080", errorHandler)
    }

    func userHandler(rw http.ResponseWriter, req *http.Request) {
        limit := GetQueries(req).Get("limit")[0]
        //...
    }
```
## More documentation comming soon