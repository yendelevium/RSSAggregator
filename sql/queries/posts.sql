-- name: CreatePost :one
INSERT INTO posts(id,
    created_at,
    update_at,
    title, 
    description,
    published_at,
    url,
    feed_id
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- Ok, this query is gonna be a little more complex
-- Basically, we just wanna get the posts from the feeds that the user is following
-- To know that, we gotta use a join, to get only the posts, who have feed_ids that the user is following
-- We also take the user_id as input as we gotta know who we want to get the feeds for
-- And also we r ordering them as most recent, and limiting how many posts we get per request

-- name: GetPostsForUser :many
SELECT posts.* from posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;