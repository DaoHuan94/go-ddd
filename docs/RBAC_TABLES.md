# RBAC Tables (Schema)

This document describes the authorization (RBAC) tables added in the RBAC migrations (`migrations/000012`..`migrations/000016`).

## Why this schema

`workspace_members` stores membership only (which users are in which workspace).

Authorization is handled by:
- `roles` (global role catalog)
- `permissions` (global permission catalog)
- `role_permissions` (role -> permissions mapping)
- `user_roles` (assign roles to a user within a workspace)

## Tables

### `roles`
Role definitions (global).

Key columns:
- `id` (PK)
- `name` (role name, unique)
- `description`
- `created_at`, `updated_at`

### `permissions`
Permission catalog.

Key columns:
- `id` (PK)
- `name` (for example: `task:read`, `task:update`, `task:delete`)
- `description`
- `created_at`, `updated_at`

### `role_permissions`
Join table mapping roles to permissions.

Key columns:
- `role_id` (FK -> `roles(id)`, `ON DELETE CASCADE`)
- `permission_id` (FK -> `permissions(id)`, `ON DELETE CASCADE`)
- composite primary key `(role_id, permission_id)`

### `user_roles`
Join table mapping a user to roles within a specific workspace.

Key columns:
- `workspace_id` (part of PK)
- `user_id` (part of PK)
- `role_id` (part of PK)
- `created_at`

Relationships:
- `(workspace_id, user_id)` references `workspace_members(workspace_id, user_id)` with `ON DELETE CASCADE`
- `role_id` references `roles(id)` with `ON DELETE CASCADE`

## Linking to existing membership

There is no role data in `workspace_members` anymore.

Instead, create rows in `user_roles` to grant roles to users in a workspace.

## Notes for application authorization

This schema defines *how* to store authorization data, but not *which* exact permissions are used.

Your application should decide a permission convention (for example: storing `resource:action` as the `permissions.name`) and then:
- insert the needed `permissions`
- assign them to each role via `role_permissions`

