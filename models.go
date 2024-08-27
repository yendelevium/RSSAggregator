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
}

// this converts sqlc User to our User, basically just copy-pasting the data into our User
func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdateAt:  dbUser.UpdateAt,
		Name:      dbUser.Name,
	}
}
