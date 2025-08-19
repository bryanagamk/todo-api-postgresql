package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"github.com/bryanagamk/todo-app-postgresql/internal/config"
	"github.com/bryanagamk/todo-app-postgresql/internal/events"
	httpx "github.com/bryanagamk/todo-app-postgresql/internal/http"
	"github.com/bryanagamk/todo-app-postgresql/internal/repo/memory"
	"github.com/bryanagamk/todo-app-postgresql/internal/worker"
)

func main() {
	cfg := config.MustLoad()

	// 1) Init App + Infra lokal
	app := fiber.New(fiber.Config{AppName: "todo-api"})
	httpx.RegisterRoutes(app) // /health

	// 2) In-memory repo (Day 2)
	repo := memory.NewTodoRepo()

	// 3) Event bus + worker pool
	bus := events.NewBus(256)
	logPool := worker.NewLoggerPool(3, "logs/todo_events.log")

	// Context untuk workers (dibatalkan saat shutdown)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	logPool.Run(workerCtx, bus)

	// 4) Register handler yang publish event
	todoHandler := httpx.NewTodoHandler(repo, bus)
	todoHandler.Register(app)

	// 5) Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	go func() {
		log.Printf("starting server on %s ...", addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	// Stop menerima event baru
	workerCancel()

	// Tutup server http dengan timeout
	_ = app.Shutdown()

	// Tutup bus (tunggu semua consumer Done)
	bus.Close()

	// Tutup resources worker
	_ = logPool.Close()

	log.Println("bye")
}
