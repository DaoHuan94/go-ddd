CREATE TABLE tags (
  id BIGSERIAL PRIMARY KEY,
  workspace_id BIGINT,
  name VARCHAR(100),

  UNIQUE (workspace_id, name),
  FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);