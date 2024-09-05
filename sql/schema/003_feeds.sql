-- Now we hv got out users set up and our authentication setup, so lets 
-- get down to some business logic. This is an RSS Feed Aggregator, so lets 
-- make a schema to store our feeds

-- The feed also has an id, created nd updated at, and a name
-- What's unique abt a feed is that it has an url, and a user_id, which is a uuid which references
-- the users id, who created the feed. Basically a foreign key to the users table
-- This also makes it such that u can't hav a feed for a user that doesn't exist

-- We are also gonna "on delete cascade", which basically means, if u delete a userid,
-- All the feeds associated with that userid will also be deleted automatically

-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;