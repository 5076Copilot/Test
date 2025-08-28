package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"mbsmf/internal/config"
	"mbsmf/internal/metrics"
	"mbsmf/internal/models"
)

type Dispatcher struct {
	cfg     config.Config
	logger  *log.Logger
	client  *http.Client
	metrics *metrics.Registry
}

func NewDispatcher(cfg config.Config, logger *log.Logger, metrics *metrics.Registry) *Dispatcher {
	return &Dispatcher{
		cfg:     cfg,
		logger:  logger,
		metrics: metrics,
		client: &http.Client{
			Timeout: cfg.HTTPClientTO,
		},
	}
}

func (d *Dispatcher) DispatchSubscriptionEvent(sub models.Subscription, event string, payload any) {
	go func() {
		body := map[string]any{
			"event":   event,
			"payload": payload,
		}
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			d.logger.Printf("notify encode error: %v", err)
			d.metrics.NotificationsFailed.Add(1)
			return
		}
		req, err := http.NewRequest(http.MethodPost, sub.CallbackURI, buf)
		if err != nil {
			d.logger.Printf("notify request create error: %v", err)
			d.metrics.NotificationsFailed.Add(1)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		start := time.Now()
		resp, err := d.client.Do(req)
		if err != nil {
			d.logger.Printf("notify send error: %v", err)
			d.metrics.NotificationsFailed.Add(1)
			return
		}
		_ = resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			d.logger.Printf("notify non-2xx: %d", resp.StatusCode)
			d.metrics.NotificationsFailed.Add(1)
			return
		}
		_ = start
		d.metrics.NotificationsSent.Add(1)
	}()
}

