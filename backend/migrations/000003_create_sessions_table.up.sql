CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INT,
  token_hash TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);