-- This folder is where we will store all our table definitions and migrations

-- goose is a database migration library- basically, git but for our database structure
-- The way goose works is that it runs the migrations in order, so the 001 table will be executed first
-- A migration is any script that will make a change to our database structure

-- Database migrations help us keep track of the STRUCTURE of the database, like creating,altering or dropping tables, 
-- adding/removing/altering columns, doing stuff with constraints and keys etc.

-- To run a goose command, u have to write 
-- goose postgres <ur DB_URL> <cmd> - we r writing postgres to tell goose its a postgres db, this is called a driver
-- goose postgres postgres://bhavikabardia:Mendelevium@localhost:5432/rssagg up - will run the up migration
-- Just run the goose help command to understand, u can set the driver and the dbstring in ur .env file so u don't always have
-- to type the postgres postgres://bhavikabardia:Mendelevium@localhost:5432/rssagg, u can just write goose up or goose status etc

-- Every migration file must have this format: 

-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    update_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);


-- +goose Down
DROP TABLE users;

-- after the goose Up comment, write ur up migration
-- down migration after ur goose Down