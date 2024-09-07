-- We r gonna need to create a new table, to store all the posts we r fetching from
-- all the different RSS Feeds

-- This is just stuff that posts have, Like title, url, the feed_id from which the post is from
-- Blah blah. U can check the RSSItem and see what all other things r there
-- Constraints are also based on that, basically just vibes
-- We r gonna make the url unique, since we don't wanna save the same post twice
-- Question: What will happen if the post is edited? Since we aren't saving the same url twice, idt it will update when we scrape
-- Maybe I can add a logic to avoid that conflict

-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;