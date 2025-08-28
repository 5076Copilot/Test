package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"example.com/mb-smf/internal/config"
	"example.com/mb-smf/internal/sbi"
)

func main() {
	cfg := config.Load()
	server := sbi.NewServer(cfg)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	if err := server.Stop(); err != nil {
		log.Printf("graceful stop error: %v", err)
	}
	log.Printf("mb-smf stopped")
}
