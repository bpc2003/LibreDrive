CREATE TABLE IF NOT EXISTS Users (
  id INTEGER PRIMARY KEY,
  username text NOT NULL UNIQUE,
  email text NOT NULL UNIQUE,
  password text NOT NULL,
  salt text NOT NULL,
  isAdmin BOOLEAN NOT NULL,
  active BOOLEAN NOT NULL
);
