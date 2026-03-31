# Project Flow (from Database Migrations)

This document explains how the domain model evolves across the SQL migrations in `./migrations/`, and how the resulting tables relate to each other.

## Migration timeline (schema evolution)

### 000001 - `users`
Creates the base identity table.

Key fields:
- `email` is unique and required (`UNIQUE`, `NOT NULL`)
- `password_hash` stored as `TEXT` (nullable in the migration)
- profile + timestamps: `name`, `avatar_url`, `created_at`, `updated_at`

### 000002 - `workspaces`
Adds organizations/containers owned by a user.

Key fields/relations:
- `workspaces.owner_id` references `users(id)`
- `created_at`, `updated_at`

### 000003 - `workspace_members`
Adds membership within a workspace.

Key fields/relations:
- `workspace_members.user_id` references `users(id)` (`ON DELETE CASCADE`)
- `workspace_members.workspace_id` references `workspaces(id)` (`ON DELETE CASCADE`)
- `workspace_members` stores only membership; RBAC assignments are stored in `user_roles`
- `UNIQUE (workspace_id, user_id)` prevents duplicate membership rows

### 000004 - `projects`
Adds projects under a workspace.

Key fields/relations:
- `projects.workspace_id` references `workspaces(id)` (`ON DELETE CASCADE`)
- `projects.created_by` references `users(id)`
- `description`, `created_at`, `updated_at`

### 000005 - `tasks`
Adds tasks under projects, with optional parent-child structure.

Key fields:
- `tasks.project_id` (required)
- `tasks.parent_task_id` (nullable) for hierarchical tasks/subtasks
- `title` required
- `status`, `priority`
- assignment and authorship: `assigned_to` (nullable) and `created_by`
- scheduling: `due_date`, `start_date`
- ordering: `position` (`INT`)
- `created_at`, `updated_at`

Key relations:
- `tasks.project_id` -> `projects(id)` (`ON DELETE CASCADE`)
- `tasks.parent_task_id` -> `tasks(id)` (`ON DELETE CASCADE`)
- `tasks.assigned_to` -> `users(id)`
- `tasks.created_by` -> `users(id)`

### 000006 - `tags`
Adds tag definitions scoped to a workspace.

Key fields/relations:
- `tags.workspace_id` references `workspaces(id)` (`ON DELETE CASCADE`)
- `UNIQUE (workspace_id, name)` ensures tag names are unique per workspace

### 000007 - `task_tags`
Adds a many-to-many relationship between tasks and tags.

Key fields/relations:
- composite primary key: `PRIMARY KEY (task_id, tag_id)`
- `task_tags.task_id` -> `tasks(id)` (`ON DELETE CASCADE`)
- `task_tags.tag_id` -> `tags(id)` (`ON DELETE CASCADE`)

### 000008 - `comments`
Adds task-related discussion.

Key fields/relations:
- `comments.task_id` -> `tasks(id)` (`ON DELETE CASCADE`)
- `comments.user_id` -> `users(id)`
- content stored as `TEXT`
- `created_at`

### 000009 - `reminders`
Adds scheduled reminders for tasks.

Key fields/relations:
- `reminders.task_id` -> `tasks(id)` (`ON DELETE CASCADE`)
- `remind_at` timestamp
- `type` categorizes the reminder (stored as `VARCHAR(50)`)
- `created_at`

### 000010 - `activity_logs`
Adds an audit/activity stream for task events.

Key fields/relations:
- `activity_logs.task_id` -> `tasks(id)` (`ON DELETE CASCADE`)
- `activity_logs.user_id` -> `users(id)`
- `action` stored as `VARCHAR(100)`
- `metadata` stored as `JSONB` for event details
- `created_at`

### 000011 - indexes
Adds performance indexes for typical query filters/joins.

Indexes created:
- `idx_tasks_project_id` on `tasks(project_id)`
- `idx_tasks_assigned_to` on `tasks(assigned_to)`
- `idx_tasks_status` on `tasks(status)`
- `idx_tasks_due_date` on `tasks(due_date)`
- `idx_comments_task_id` on `comments(task_id)`
- `idx_activity_task_id` on `activity_logs(task_id)`

## Resulting entity relationship "flow" (final schema)

After `000011`, the core flow of data/entities is:

1. `users`
2. `workspaces` owned by a `user` (`workspaces.owner_id`)
3. `workspace_members` associates additional `users` to a `workspace`
4. `projects` belong to a `workspace` (`projects.workspace_id`)
5. `tasks` belong to a `project` (`tasks.project_id`)
6. `tasks` optionally reference other `tasks` (`tasks.parent_task_id`) for subtasks/hierarchy
7. `tags` are defined per `workspace`
8. `tasks` and `tags` are linked via `task_tags` (many-to-many)
9. `comments`, `reminders`, and `activity_logs` hang off `tasks`

## Delete semantics (cascade behavior)

The migrations encode a consistent "container deletes children" pattern via `ON DELETE CASCADE`:
- Deleting a `workspace` cascades to:
  - `workspace_members`
  - `projects`
  - `tags`
  - `user_roles`
- Deleting a `project` cascades to `tasks`
- Deleting a `task` cascades to:
  - `task_tags`
  - `comments`
  - `reminders`
  - `activity_logs`
  - and also cascades to subtasks through `tasks.parent_task_id`

## How to read the migration "flow" in development

- Treat `000001` -> `000004` as the "identity + container" foundation.
- Treat `000005` as the central entity: almost everything else attaches to `tasks`.
- Treat `000006` + `000007` as "classification" (tags) layered on top of tasks.
- Treat `000008` -> `000010` as "task context" (comments, reminders, audit/activity).
- Treat `000011` as "query optimization" after the first functional set of relations is in place.

`users -> workspaces -> projects -> tasks ties to tags, comments, reminders, and activity_logs`