-- Let a user follow a feed
-- name: CreateFeedFollow :one
INSERT INTO feed_follows(id,created_at,update_at,user_id,feed_id)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;


-- Let's get a way for the user to see all the feeds he's following
-- name: GetFeedFollows :many
SELECT * FROM feed_follows WHERE user_id=$1;

-- We also need a way to unfollow feeds, which is basically deleting the record in the feed_follows table
-- This is gonna be a query that does't return anything, it will just execute. Hence the :exec
-- We don't really need the user_id for this to work, as the id already identifies the feed_follow
-- But this will prevent someone other than the user from deleting the feed follow of that user

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id=$1 AND user_id=$2;