CREATE TABLE IF NOT EXISTS events(
id SERIAL PRIMARY KEY,
name VARCHAR(100) NOT NULL,
description TEXT,
location TEXT,
address TEXT,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
event_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rsvps(
id SERIAL PRIMARY KEY,
event_id integer REFERENCES events,
guest_name VARCHAR(100) NOT NULL,
email VARCHAR(100) NOT NULL,
number_guests NUMERIC,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);