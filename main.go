package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"findit-campus-backend/app/repository"
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/database"
	"findit-campus-backend/route"
)

func main() {
	cfg := config.LoadConfig()
	logger := config.InitLogger(cfg.AppEnv)

	dbPool, err := database.NewPostgresPool(cfg)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Dependency Injection (dari dalam ke luar)
	userRepo := repository.NewUserRepository(dbPool)
	authService := service.NewAuthService(userRepo, cfg)
	returnHistoryRepo := repository.NewReturnHistoryRepository(dbPool)
	returnHistoryService := service.NewReturnHistoryService(returnHistoryRepo)

	reportRepo := repository.NewReportRepository(dbPool)
	reportService := service.NewReportService(reportRepo, cfg)

	app := config.NewApp(cfg, logger)
	route.Register(app, cfg, route.Services{
		Auth:   authService,
		Report: reportService,
		ReturnHistory: returnHistoryService,
	})

	// Jalankan server di goroutine dengan graceful shutdown
	go func() {
		addr := ":" + cfg.AppPort
		logger.Info("server starting", "addr", addr, "env", cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			logger.Error("server stopped unexpectedly", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}
	logger.Info("server exited")
}