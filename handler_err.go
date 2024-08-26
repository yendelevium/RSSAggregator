package main

import "net/http"

// The handler for errors
// We r responding with the same error everytime here
// 400 status code coz client error
func handlerErr(w http.ResponseWriter, r *http.Request) {
	repsondWithError(w, 400, "Something went wrong")
}
