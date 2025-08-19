package domain_test

import (
	"testing"
	"time"

	"github.com/bryanagamk/todo-app-postgresql/internal/domain"
)

func TestNewTodo_Valid(t *testing.T) {
	td, err := domain.NewTodo("Belajar Go", "Day 1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if td.Status != domain.StatusPending {
		t.Fatalf("want status pending, got %s", td.Status)
	}
	if td.Title != "Belajar Go" {
		t.Fatalf("wrong title")
	}
}

func TestNewTodo_InvalidTitle(t *testing.T) {
	_, err := domain.NewTodo("  a ", "", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != domain.ErrInvalidTitle {
		t.Fatalf("want ErrInvalidTitle, got %v", err)
	}
}

func TestNewTodo_InvalidDue(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	_, err := domain.NewTodo("Belajar", "", &past)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != domain.ErrInvalidDueBefore {
		t.Fatalf("want ErrInvalidDueBefore, got %v", err)
	}
}

func TestTodo_Validate_Status(t *testing.T) {
	td, _ := domain.NewTodo("Belajar", "", nil)
	td.Status = "weird"
	if err := td.Validate(); err == nil {
		t.Fatalf("expected invalid status error")
	}
}

func TestTodo_MarkDone(t *testing.T) {
	td, _ := domain.NewTodo("Belajar", "", nil)
	td.MarkDone()
	if td.Status != domain.StatusDone {
		t.Fatalf("want done, got %s", td.Status)
	}
}
