CREATE TABLE projects (
  id BIGSERIAL PRIMARY KEY,
  workspace_id BIGINT,
  name VARCHAR(255),
  description TEXT,
  created_by BIGINT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id)
);