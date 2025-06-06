package main

import (
	"github.com/Dor1ma/Scheduler/internal/api"
	"github.com/Dor1ma/Scheduler/internal/app"
	"github.com/Dor1ma/Scheduler/internal/config"
	"github.com/Dor1ma/Scheduler/internal/storage/postgres"
	"github.com/Dor1ma/Scheduler/internal/timer"
	"github.com/Dor1ma/Scheduler/pkg/logger"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	cfg := config.Load()

	store, err := postgres.NewPostgresStore(cfg.GetPostgresDSN())
	if err != nil {
		logger.Fatal("failed to connect to DB: %v", err)
	}

	dispatcher := app.NewDispatcher(10000)

	for i := 0; i < 10; i++ {
		worker := app.NewWorker(i, dispatcher)
		worker.Start()
	}

	wheel := timer.NewTimeWheel(1*time.Second, 3600, func(task timer.Task) {
		dispatcher.Dispatch(task)
	})

	scheduler := app.NewScheduler(wheel, store)
	if err := scheduler.Start(); err != nil {
		logger.Fatal("scheduler start failed: %v", err)
	}

	fiberApp := fiber.New()
	api.RegisterRoutes(fiberApp, scheduler)

	logger.Info("Scheduler API listening on :%s", cfg.AppPort)
	if err := fiberApp.Listen(":" + cfg.AppPort); err != nil {
		logger.Fatal("server error: %v", err)
	}
}
