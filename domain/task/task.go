package task

// Task is the domain representation of a task.
// It is used by usecases and repository interfaces.
type Task struct {
	ID          int64
	ProjectID   int64
	Title       string
	Description string

	Priority string
}

type CreateTaskInput struct {
	ProjectID   int64
	Title       string
	Description string
	Priority    string
}

