-- We've built out and completed most of the CRUD part of our API, but we haven't
-- Built the most important part, the part of the server that actually goes
-- And fetched out the posts from the RSSFeeds stored in our db
-- The main purpose of this server, is so that it can go out periodically, and download
-- All of the posts on the rssFeeds. This way it keeps checking if there's a new post, and if there 
-- is, we download the post and show it to the user

-- So, to do this, we first have to add a new column to our feeds table, called last_fetched_at
-- So we know when was the last time the scraper went and fetched the posts for a given feed
-- We are letting it be NULL, incase the feed has never been fetched
-- Since we are defaulting it to be null, we don't need to update the CreateFeed query

-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;