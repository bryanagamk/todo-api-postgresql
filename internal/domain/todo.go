package domain

import (
	"errors"
	"strings"
	"time"
)

type TodoStatus string

const (
	StatusPending TodoStatus = "pending"
	StatusDone    TodoStatus = "done"
)

var (
	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrInvalidDueBefore = errors.New("invalid due: due date before created at")
)

type Todo struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Note      string     `json:"note,omitempty"`
	Status    TodoStatus `json:"status"`
	DueAt     *time.Time `json:"due_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewTodo(title, note string, dueAt *time.Time) (*Todo, error) {
	title = strings.TrimSpace(title)
	if len(title) < 3 {
		return nil, ErrInvalidTitle
	}
	now := time.Now().UTC()

	t := &Todo{
		Title:     title,
		Note:      strings.TrimSpace(note),
		Status:    StatusPending,
		DueAt:     dueAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if dueAt != nil && dueAt.Before(now) {
		return nil, ErrInvalidDueBefore
	}
	return t, nil
}

func (t *Todo) Validate() error {
	if strings.TrimSpace(t.Title) == "" || len(t.Title) < 3 {
		return ErrInvalidTitle
	}
	switch t.Status {
	case StatusPending, StatusDone:
	default:
		return ErrInvalidStatus
	}
	if t.DueAt != nil && t.DueAt.Before(t.CreatedAt) {
		return ErrInvalidDueBefore
	}
	return nil
}

func (t *Todo) MarkDone() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now().UTC()
}
