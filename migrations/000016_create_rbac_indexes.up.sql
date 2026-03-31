-- RBAC: indexes for common joins/filters

CREATE INDEX IF NOT EXISTS idx_user_roles_workspace_id ON user_roles(workspace_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

