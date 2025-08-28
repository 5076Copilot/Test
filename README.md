## MB-SMF (Multicast-Broadcast SMF) - Minimal Reference Node

This repository provides a minimal MB-SMF node skeleton in Go. It exposes a basic Nmbsf v1-like API surface for sessions, subscriptions, and policies, health and Prometheus-style metrics, a stubbed NRF heartbeat, and an in-memory store.

Note: This is a reference skeleton for experimentation and not a full 3GPP-compliant implementation.

### Build

```bash
make build
```

### Run

```bash
make run
```

Environment variables:

- MB_SMF_HTTP_ADDR: HTTP listen address (default ":8080")
- MB_SMF_NRF_ENABLED: enable NRF heartbeat (default false)
- MB_SMF_NRF_URI: NRF URI (required when enabled)
- MB_SMF_NODE_ID: node identifier (default "mbsmf-1")

### API

- GET /healthz
- GET /metrics

- POST /nmbsf/v1/sessions
- GET  /nmbsf/v1/sessions
- DELETE /nmbsf/v1/sessions/{id}

- POST /nmbsf/v1/subscriptions
- GET  /nmbsf/v1/subscriptions
- DELETE /nmbsf/v1/subscriptions/{id}

- POST /nmbsf/v1/policies
- GET  /nmbsf/v1/policies
- DELETE /nmbsf/v1/policies/{id}

# Test