-- name: CreateUser :one
INSERT INTO users(id,created_at,update_at,name, api_key)
VALUES ($1,$2,$3,$4, encode(sha256(random()::text::bytea),'hex') )
RETURNING *;

-- The way sqlcn works is that it takes the sql query, and creates type-safe go code which matches the query
-- The way sqlcn works is that it takes the sql query, and creates type-safe go code which matches the query
-- Every sqlc query starts with an sql comment, name: <queryname> :<no.of records to be returned by this query>
-- wtf are the $thingys u ask? 
-- In sqlc, each $number represent parameters for the function. This statement creates a function, which takes 
-- 4 arguments, and the first argument is put in place of $1 and so on
-- RETURNING *; is that we r creating a new record, and we wanna return that record

-- We always run sqlc from the root of our package, where the sqlc.yaml file is located
-- we write sqlc generate in the cmd
-- So what this does, is that sql has access to our schema and all queries, 
-- as specified in the sqlc.yaml file, and it goes and generates the go code in the internal/database repo,
-- which was again specified in the sqlc.yaml file

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;