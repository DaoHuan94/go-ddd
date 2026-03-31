CREATE TABLE workspaces (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255),
  owner_id BIGINT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  FOREIGN KEY (owner_id) REFERENCES users(id)
);