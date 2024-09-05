-- We've given users a way to create feed and a way to get all the feeds
-- Now lets give the users a way to follow specific feeds
-- This table will just store a relationship b/w a user and all of the feeds they follow
-- So it's gonna be a many to many relationship

-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id,feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
-- The unique constraint makes it so that we can never have 2 instances of a follow
-- For the same user and feed. A unique user can only follow a certain feed once, u can't really follow it twice rite
