package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/bryanagamk/todo-app-postgresql/internal/domain"
	"github.com/bryanagamk/todo-app-postgresql/internal/events"
)

type TodoCreator interface {
	Create(ctx context.Context, t *domain.Todo) (*domain.Todo, error)
}

type TodoHandler struct {
	repo TodoCreator
	bus  *events.Bus
}

func NewTodoHandler(repo TodoCreator, bus *events.Bus) *TodoHandler {
	return &TodoHandler{repo: repo, bus: bus}
}

type createTodoReq struct {
	Title string     `json:"title"`
	Note  string     `json:"note"`
	DueAt *time.Time `json:"due_at"`
}

func (h *TodoHandler) Register(app *fiber.App) {
	app.Post("/todos", h.CreateTodo)
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req createTodoReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	td, err := domain.NewTodo(req.Title, req.Note, req.DueAt)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Gunakan context dengan timeout agar tidak menggantung
	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	created, err := h.repo.Create(ctx, td)
	if err != nil {
		code := http.StatusInternalServerError
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}

	// PUBLISH EVENT (non-blocking, hormati ctx)
	_ = h.bus.Publish(ctx, events.NewTodoCreated(created.ID, created.Title))

	return c.Status(http.StatusCreated).JSON(created)
}
