package nrf

import (
	"context"
	"log"
	"time"

	"mbsmf/internal/config"
)

type Client struct {
	cfg    config.Config
	logger *log.Logger
}

func NewClient(cfg config.Config, logger *log.Logger) *Client {
	return &Client{cfg: cfg, logger: logger}
}

// StartBackground simulates periodic registration/heartbeat with NRF.
func (c *Client) StartBackground(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if c.cfg.NRFURI == "" {
			c.logger.Println("NRF client enabled but NRF_URI is empty; skipping")
			return
		}
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		c.logger.Printf("NRF heartbeat started for node %s -> %s", c.cfg.NodeID, c.cfg.NRFURI)
		for {
			select {
			case <-ctx.Done():
				c.logger.Println("NRF heartbeat stopped")
				return
			case <-ticker.C:
				c.logger.Printf("NRF heartbeat: node=%s", c.cfg.NodeID)
			}
		}
	}()
	return cancel
}

