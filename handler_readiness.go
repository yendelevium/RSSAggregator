// json.go helps us to respond with json data
// So now, lets create an HTTP handler that responds with the json data(basically a function which will call that function)
package main

import "net/http"

// This is a very specific function signature, a function signature you HAVE to use
// if u wanna define a HTTP handler in the way the go standard library expects
// U take a responseWriter as the first param, and a POINTER to the http request as the 2nd
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Rn, we only care about the status code, so in the payload let's just send an empty struct, which will marshal to an empty json string
	respondWithJSON(w, 200, struct{}{})
}
