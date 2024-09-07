-- We gonna create a query to create a new feed, and to get a new feed

-- name: CreateFeed :one
INSERT INTO feeds(id,created_at,update_at,name,url,user_id)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- This query is to get all the feeds from our db
-- Hence we use :many as many records can be returned

-- name: GetFeeds :many
SELECT * FROM feeds;


-- This function will go get the feed, that next needs to be fetched
-- First, we wanna find feeds that have never been fectched, and then ordering them, by most recently fetched/ most unrecently fetched, idk how dates work in sql
-- We are also asking the user how many feeds they want

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- This is the one we call after we fetch the feed,to update it,and return the updated feed
-- The updated_at and created_at fields are mostly for auditing purposes.
-- it's pretty standard practice to set these on every sql record, to see when they were created nd updated

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
update_at = NOW()
WHERE id = $1
RETURNING *;