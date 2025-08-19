package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/bryanagamk/todo-app-postgresql/internal/config"
	httpx "github.com/bryanagamk/todo-app-postgresql/internal/http"
)

func main() {
	cfg := config.MustLoad()

	app := fiber.New(fiber.Config{
		AppName: "todo-api",
	})

	httpx.RegisterRoutes(app)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s ...", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
