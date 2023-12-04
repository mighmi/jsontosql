CREATE TABLE people (
  id SERIAL PRIMARY KEY,
  FirstName TEXT,
  LastName TEXT,
  Latitude DOUBLE PRECISION,
  Longitude DOUBLE PRECISION,
  Username TEXT,
  passwd TEXT, 
  Email TEXT UNIQUE NOT NULL,
  DateOfBirth Date
);
