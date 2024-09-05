-- We gonna create a query to create a new feed, and to get a new feed

-- name: CreateFeed :one
INSERT INTO feeds(id,created_at,update_at,name,url,user_id)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- This query is to get all the feeds from our db
-- Hence we use :many as many records can be returned
-- name: GetFeeds :many
SELECT * FROM feeds;