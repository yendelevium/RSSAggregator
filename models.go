package main

import (
	"time"

	"github.com/Yendelevium/RSSAggregator/internal/database"
	"github.com/google/uuid"
)

// Just creating our own User struct, identical to the one created by sqlc
// We r just adding json-reflect tags to get it how we want in the json response
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

// this converts sqlc User to our User, basically just copy-pasting the data into our User
func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdateAt:  dbUser.UpdateAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedtoFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdateAt:  dbFeed.UpdateAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func databaseFeedstoFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		// Converting every feed in the database.Feed slice to OUR feed type
		feeds = append(feeds, databaseFeedtoFeed(dbFeed))
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowtoFeedFollow(dbFeed database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdateAt:  dbFeed.UpdateAt,
		UserID:    dbFeed.UserID,
		FeedID:    dbFeed.FeedID,
	}
}

func databaseFeedFollowstoFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}
	for _, dbFeedFollow := range dbFeedFollows {
		// Converting every feed in the database.Feed slice to OUR feed type
		feedFollows = append(feedFollows, databaseFeedFollowtoFeedFollow(dbFeedFollow))
	}
	return feedFollows
}

// We don't want the sql.NullString here, as this struct will be marshalled as a json
// Since sql.NullString is a struct, when we marshall, we will get "description", "string", and "valid"
// as json keys. We only need the description, no need for these nested objects
// Nested json stuff is pretty bad user experience, considering that JSON supports NULL as a value
// So if the description is empty, we can just show NULL
// So we will make a pointer to a string. This is coz the way that json works, is that
// if we marshall a pointer to a string, and its null, json will marhsall that to null in json, otherwise the value of the pointer
// U can't just have an empty string fr desc, as then that means the description exists, just nothing is in it
// Null means the description XML tag is not present itself in the RSSFeed
type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdateAt:    dbPost.UpdateAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	posts := []Post{}
	for _, dbPost := range dbPosts {
		posts = append(posts, databasePostToPost(dbPost))
	}
	return posts
}
