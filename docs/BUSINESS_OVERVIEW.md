# Business Overview

This project is a multi-tenant task management system organized around the idea of **workspaces**.

## Core concepts

### Users
An account identified by `users.id` (email + profile fields).

### Workspaces
A workspace is a top-level container owned by a user (`workspaces.owner_id`).
Workspaces have members (`workspace_members`) which control who can access workspace data.

### Projects
A project belongs to a workspace (`projects.workspace_id`) and groups tasks.

### Tasks
Tasks belong to projects (`tasks.project_id`).
Tasks support:
- hierarchical subtasks via `tasks.parent_task_id`
- ordering via `tasks.position`
- workflow metadata such as `status`, `priority`, `due_date`, `start_date`
- assignment to users via `tasks.assigned_to`

### Tags
Tags are defined per workspace (`tags.workspace_id`) and can be attached to tasks via `task_tags` (many-to-many).

### Comments
Comments are attached to tasks (`comments.task_id`) and authored by users (`comments.user_id`).

### Reminders
Reminders are attached to tasks (`reminders.task_id`) with a `remind_at` timestamp and a `type`.

### Activity Logs
Activity logs are an audit/event stream for tasks (`activity_logs.task_id`) with:
- `action` (string)
- `metadata` (JSONB payload)
- author/user (`activity_logs.user_id`)

## Authorization (RBAC)

Authorization is workspace-scoped and modeled as:
- `roles` (role catalog)
- `permissions` (permission catalog; commonly modeled as `resource:action`)
- `role_permissions` (role -> permissions)
- `user_roles` (assignment of roles to users within a workspace)

At runtime, the effective permissions for a user in a workspace come from:
`workspace_members` (membership) -> `user_roles` (role assignments) -> `role_permissions` -> `permissions`.

