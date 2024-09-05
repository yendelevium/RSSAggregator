package main

import (
	"fmt"
	"net/http"

	"github.com/Yendelevium/RSSAggregator/internal/auth"
	"github.com/Yendelevium/RSSAggregator/internal/database"
)

// Here, we are gonna create a new type, which looks ALMOST like a regular http handler function,
// But it includes a third parameter, which is the user associated with it
// If u think about it, any authenticated handler, will have the authenticated user associated with it

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// The problem with this authHandler type we created, is that it doesn't match the function  signature
// of an http.HandlerFunc

// So we are gonna create a new function, which willl be a method on the apiCfg so that it can access the db
// And this function will take an authedHandeler, and will return a http.HandlerFunc, so that we can use it with the chi router

// They way this function will actually work, is that we are gonna return a CLOSURE
// This way we can access the authed-handler inside the http.HandlerFunc we are returning

// We'll return an anonymous function, which will have the same function as an http.HandlerFunc
// So we will just rip out the code from the handlerGetUser function, and put it here
// Then we can just call the authedHandler, and we can pass in the user, and do the handler-specific stuff
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We have the GetAPIKey function, so now we can use it
		// We pass the http headers, which is given by r.Header
		apiKey, err := auth.GetAPIKey(r.Header)
		// If there's an error while getting the APIKey, respond with an error
		if err != nil {
			// 403 is an error code for like a permission error
			repsondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		// Now that we have our apiKey, we can use our db query, to get the user by APIKey
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			repsondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		// Since this is a closure, we can access the authed-handler here
		// We can pass the user to the handler, and then let the hander do its thing
		handler(w, r, user)
	}
}
