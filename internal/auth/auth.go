package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from the headers of an HTTP request
// Basically it will look for a specific (http)header and, if it finds the apikey,
// It will return it, otherwise it wll return an error

// As the authors of this server, we get to decide how the header will look like
// And by "We" i really mean him coz he's the one showing me how to code this

// Example: key : value
// Authorization : ApiKey <the apikey>

func GetAPIKey(headers http.Header) (string, error) {
	// headers.Get(<header>) returns that HTTP Header's value
	val := headers.Get("Authorization")

	// If no header is found, raise an error
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	// strings.Split(str,<delimiter>) splits the string like in python on that delimitort,
	// returns a slice of strings like in python again
	// We split it on spaces, as we are expecting the key to be like "ApiKey <the apikey>"
	vals := strings.Split(val, " ")

	// If the key isnt like how we want, raise error
	if len(vals) != 2 {
		// also ur not supposed to have the first letter as a capital letter in an error in go
		return "", errors.New("malformed auth header")
	}

	// The first word isnt ApiKey, raise error
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}

	// If all is good, return the APIKey and nil error
	return vals[1], nil
}
