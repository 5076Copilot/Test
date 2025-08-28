package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mbsmf/internal/config"
	"mbsmf/internal/httpserver"
	"mbsmf/internal/metrics"
	"mbsmf/internal/nrf"
	"mbsmf/internal/notify"
	"mbsmf/internal/service"
	memstore "mbsmf/internal/store/memory"
)

func main() {
	cfg := config.LoadFromEnv()

	logger := log.New(os.Stdout, "mbsmf ", log.LstdFlags|log.Lmicroseconds|log.LUTC)

	store := memstore.NewMemoryStore()
	metric := metrics.NewRegistry()
	notifier := notify.NewDispatcher(cfg, logger, metric)
	svc := service.NewMBService(store, notifier, logger, metric, cfg)

	srv := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           httpserver.NewRouter(svc, logger, metric, cfg),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start NRF stub if enabled
	var nrfStop func()
	if cfg.NRFEnabled {
		client := nrf.NewClient(cfg, logger)
		nrfStop = client.StartBackground(ctx)
	}

	go func() {
		logger.Printf("HTTP server listening on %s", cfg.HTTPAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Println("shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if nrfStop != nil {
		nrfStop()
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}
	logger.Println("server stopped")
}

