// The server we r building is a JSON REST API
// This means all the request bodies coming in and response bodies going back will have a JSON format
// This function will help make it easier to send the response bodies
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// This function creates a JSON response
// https://hackthedeveloper.com/golang-responsewriter-request/
// a http.ResponseWriter is an interface that provides a way for the server to construct an HTTP Response to the clients request
// code is the status code
// payload is the data u wanna respond with in the JSON. It can be any type, so we use the empty interface{}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// json.Marshal() converts a data structure, like a slice or a struct into a json string, and it returns
	// it as bytes, and the reason it returns it as bytes, is so we can write it in a binary format, directly to the HTTP Response, which is pretty convinient
	// json.Unmarshal() is the reverse process. This is called JSON Marshalling and Unmarshalling
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		// This will give a status code of 500. WriteHeader is a method on the ResponseWriter interface
		w.WriteHeader(500)
		return
	}
	// IF payload is marshalled successfully, then :
	// The Header() http.Header function gives you access to the response headers.
	// It returns an http.Header type, which you can use to manipulate headers such as Content-Type, Cache-Control, and more.
	w.Header().Add("Content-Type", "application/json")
	// We r saying that the response content is a json.
	// Adds a response header "Content-Type" with value "application/json", which is the standard value fr JSON responses
	// Sending a status code of 200, which means OK
	w.WriteHeader(code)
	// Sendind the data u wanna send. It should be a byte slice!
	w.Write(dat)
}

// This is a function to respond with if there's an error. very similar to the one above, but takes a msg string instead of a payload
// It will format the msg string into a consistent json object every single time
func repsondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		// We r just logging this, coz codes in the 500+ range are server side- errors
		// So logging this will let us know it's a mistake/bug on our side
		log.Println("Responding with 5XX error:", msg)
	}

	// This creates a specific strucure of json
	type errResponse struct {
		// The `json:"error"` is to say that the key this struct will get marshalled to should be "error"
		// These json-reflect tags in a struct are typically used to tell the json.Marshall and json.Unmarshall
		// how we want them to convert this struct into a json object
		Error string `json:"error"`
	}
	// Our json will look somethinng like
	// {
	// 	"error" : "Err Msg"
	// }

	// Now we r just calling the function to respond with the err msg as a json
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
