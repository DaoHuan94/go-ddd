# API Reference

## Base URL

`/api/v1`

## Implemented endpoints (current code)

### Get task by id
`GET /api/v1/tasks/:id`

Behavior:
- Returns a task summary (`id`, `project_id`, `title`, `description`).

Request:
- The current controller binds JSON into `GetTaskByIDRequest` (`{"id": <int>}`) using `c.Bind`, even though the route uses `:id`.
- Practically, the request expects JSON body `{"id": 123}`.

Responses:
- `200 OK`
```json
{
  "message": "success",
  "data": {
    "id": 1,
    "project_id": 10,
    "title": "Example",
    "description": "..."
  }
}
```
- `501 Not Implemented` if the usecase dependency is missing
- `400 Bad Request` if request binding fails
- `500 Internal Server Error` on unexpected repository errors

### Authentication: register/login/refresh/logout

Base: `/api/v1/auth`

1. `POST /api/v1/auth/register`
Request:
```json
{
  "email": "user@example.com",
  "password": "string",
  "name": "optional",
  "avatar_url": "optional"
}
```
Response (`200 OK`):
```json
{
  "message": "success",
  "data": {
    "access_token": "string",
    "refresh_token": "string"
  }
}
```

2. `POST /api/v1/auth/login`
Request:
```json
{
  "email": "user@example.com",
  "password": "string"
}
```
Response (`200 OK`): same `data` shape as register.

3. `POST /api/v1/auth/refresh`
Request:
```json
{
  "refresh_token": "string"
}
```
Response (`200 OK`): returns a new `access_token` + `refresh_token`.

4. `POST /api/v1/auth/logout`
Request:
```json
{
  "refresh_token": "string"
}
```
Response (`200 OK`):
```json
{ "message": "success" }
```

## Suggested REST APIs for the full domain model (not implemented yet)

The following endpoints are aligned to the tables in your migrations and follow conventional REST patterns. They are *proposals* unless you wire controllers/usecases for them.

### Workspaces
- `POST /workspaces`
- `GET /workspaces/:workspace_id`
- `POST /workspaces/:workspace_id/members`
- `POST /workspaces/:workspace_id/users/:user_id/roles/:role_id`

### Roles/Permissions (RBAC)
- `POST /workspaces/:workspace_id/roles`
- `POST /permissions`
- `POST /roles/:role_id/permissions`

### Projects
- `POST /workspaces/:workspace_id/projects`
- `GET /projects/:project_id`

### Tasks
- `POST /projects/:project_id/tasks`
- `GET /tasks/:task_id`
- `PUT /tasks/:task_id`
- `DELETE /tasks/:task_id`

### Subtasks
- create with `tasks.parent_task_id` via `POST /tasks`

### Tags and Task tags
- `POST /workspaces/:workspace_id/tags`
- `POST /tasks/:task_id/tags/:tag_id`
- `DELETE /tasks/:task_id/tags/:tag_id`

### Comments, Reminders, Activity logs
- `POST /tasks/:task_id/comments`
- `GET /tasks/:task_id/comments`
- `POST /tasks/:task_id/reminders`
- `GET /tasks/:task_id/reminders`
- `GET /tasks/:task_id/activity_logs`

