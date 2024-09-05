package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yendelevium/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// handlerCreateFeed will have the same structure like handlerCreateUser
// But, to create a feed, we need the user, who we will get by an APIKey
// So it will need a lot of same logic like the handlerGetUser function
// So, instead of copypasting the 10-15 lines of code from handlerGetUSer to every handler that gets authentocated,
// We will create a MIDDLEWARE, to kindof DRY-up the code
func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// We have access to the user thank's to the middleware
	// Now we r just creating the feed, very similar to creating the user,
	// but we have an extra parameter, which is the url of the feed, which we want from the user
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Error parsion JSON: %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedtoFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	// feeds is a slice of database.feeds, as it can be more than one row
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		repsondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}
	// This will return an array of JSON Objects
	// Make the get request and see
	respondWithJSON(w, 201, databaseFeedstoFeeds(feeds))
}
