CREATE TABLE IF NOT EXISTS Users (
  id INTEGER PRIMARY KEY,
  username text NOT NULL UNIQUE,
  password text NOT NULL,
  isAdmin BOOLEAN NOT NULL
);
