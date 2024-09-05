-- We'll be using APIKeys to authenticate our users on the server
-- An APIKey is a code used to identify an application or user and is used for authentication in computer applications.
-- Not only is an APIKey more secure than an id-password, it's so long that it can be used to uniquely identify a user
-- These APIKeys must also be kept secret by the user coz just the APIKey can be used to authenticate the user
-- Usually, once we make a migration, we don't change it, coz it can mess some data shit up
-- Hence, to add this APIKey column, we are making another migration

-- We set a default APIKey cause we alr have users in our database, and this column is unique and not null
-- Sql wont let us ad this column as the other 2 records don't have an apikey 
-- varchar(64) as we want our key to be 64 characters long
-- encode(sha256(random()::text::bytea),'hex') This generates a random apiKey 
-- We are basically generating some random bytes, and we are casting it into a byte array
-- then we are using the sha256 HASH FUNCTION to get a fixed size output, which is 64 here
-- And then encode this in hexadecimal

-- Here's a breakdown of the expression:
-- random(): Generates a random number between 0 and 1.
-- random()::text: Converts the random number to a text format.
-- random()::text::bytea: Further converts the text to a bytea (binary data) type.
-- sha256(...): Applies the SHA-256 hash function to the binary data, resulting in a 256-bit hash.
-- encode(..., 'hex'): Converts the resulting hash from binary to a hexadecimal string representation.

-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT(
    encode(sha256(random()::text::bytea),'hex')
);

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;