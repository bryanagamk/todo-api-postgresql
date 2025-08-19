package events

import "time"

type EventType string

const (
	EventTodoCreated EventType = "todo.created"
)

type TodoCreatedPayload struct {
	ID    int64
	Title string
}

type Event struct {
	Type      EventType
	Timestamp time.Time
	Payload   interface{}
}
