package mux

// Classic returns a new classic router instance.
func Classic() *Router {
	router := NewRouter()
	router.UseRoute(newRoute)
	return router
}
