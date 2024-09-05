package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yendelevium/RSSAggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// In order to make a user follow a feed, all we need to do is to create a feed_follows record
// woth that user-feed relationship
// This will also be an authenticated endpoint, as we need a user and we need them to be authenticated
// We'll also need them to tell us which feed they wanna follow, which will be the feed_id
func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowtoFeedFollow(feedFollow))
}

// This is also authenticated, as we need the userID to get the feeds hes following
func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowstoFeedFollows(feedFollows))
}

// Now, to delete a feed follow, we will need a feed_follow_id
// Delete requests, usually don't have a body in the payload. It's not that conventionally basically
// The payload of an API is the data you are interested in transporting to the server when you make an API request.
// Simply put, it is the body of your HTTP request and response message.
// It's a little more conventional to pass the id in the http-path, basically the url where u send the delete request
func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// To get the feedFollowId, the chi router has the URLParam() function, where u pass
	// in the http request, and the key, which is the thing u wrote in the {} in the http path while hooking up the handler
	// This returns a string
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	// We will parse the str to a uuid
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	// What matters to the client is simply the 200 status code, so we r just gonna repond with an empty JSON
	// But u CAN put json obj with a msg or smtg like "Deleted the feed follow" or smtg if u want
	respondWithJSON(w, 200, struct{}{})
}
