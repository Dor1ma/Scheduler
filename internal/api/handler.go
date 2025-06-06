package api

import (
	"encoding/json"
	"github.com/Dor1ma/Scheduler/internal/app"
	"github.com/Dor1ma/Scheduler/internal/timer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type createTaskRequest struct {
	ExecuteAt time.Time       `json:"execute_at"`
	Method    string          `json:"method"`
	URL       string          `json:"url,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

func RegisterRoutes(app *fiber.App, scheduler *app.Scheduler) {
	app.Post("/task", func(c *fiber.Ctx) error {
		var req createTaskRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request",
			})
		}

		id := uuid.New().String()

		task := timer.Task{
			ID:        id,
			ExecuteAt: req.ExecuteAt,
			Method:    req.Method,
			URL:       req.URL,
			Payload:   req.Payload,
		}

		if err := scheduler.ScheduleTask(task); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	})

	app.Delete("/task/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := scheduler.CancelTask(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
}
