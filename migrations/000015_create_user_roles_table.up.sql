-- RBAC: user_roles join table (user <-> role within a workspace)

CREATE TABLE IF NOT EXISTS user_roles (
  workspace_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),

  PRIMARY KEY (workspace_id, user_id, role_id),

  -- Enforce that role assignment maps to a real workspace membership.
  FOREIGN KEY (workspace_id, user_id) REFERENCES workspace_members(workspace_id, user_id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

