MB-SMF (Multicast-Broadcast SMF) – Reference Skeleton

Overview

This repository provides a minimal Go implementation skeleton of an MB-SMF for 5G Multicast-Broadcast Services (MBS). It exposes simplified SBI-style HTTP endpoints (inspired by 3GPP TS 29.532) and includes a PFCP client stub (3GPP TS 29.244) for session signaling towards the User Plane. The radio and RAN aspects are outside of scope; this focuses on control-plane MB-SMF basics.

Key references

- 3GPP TS 23.247 (MBS architecture)
- 3GPP TS 29.532 (MBS Session Management Services)
- 3GPP TS 29.537 (MBS Policy Control services)
- 3GPP TS 29.244 (PFCP)
- 3GPP TS 29.580/29.581 (MB-SMF, MBSTF)

Project layout

- cmd/mb-smf: entrypoint main
- internal/config: env-based configuration
- internal/domain: domain models (MBS session, requests)
- internal/repo: in-memory repository
- internal/service: session lifecycle orchestration
- internal/pfcp: minimal PFCP client stub (UDP, minimal header encode)
- internal/sbi: HTTP server and endpoints

Build and run

```bash
cd /workspace/mb-smf
go build ./...
MB_SMF_ADDR=":8080" MB_SMF_PFCP_ADDR="127.0.0.1:8805" go run ./cmd/mb-smf
```

Configuration

- MB_SMF_ADDR: HTTP listen address (default ":8080")
- MB_SMF_PFCP_ADDR: PFCP peer (UPF) address, e.g. "127.0.0.1:8805"

HTTP API (simplified)

- POST /mbs/sessions
  - Body: {"serviceArea":"SA1","contentDesc":"news","policies":{"qos":"mbs"}}
  - Returns 201 with session JSON
- GET /mbs/sessions
- GET /mbs/sessions/{id}
- PATCH /mbs/sessions/{id}
  - Body: {"policies":{...}, "state":"ACTIVE|INACTIVE"}
- POST /mbs/sessions/{id}?action=activate
- POST /mbs/sessions/{id}?action=deactivate
- DELETE /mbs/sessions/{id}

Example curl

```bash
# Create session
curl -sS -X POST http://localhost:8080/mbs/sessions \
  -H 'content-type: application/json' \
  -d '{"serviceArea":"SA1","contentDesc":"news","policies":{"qos":"mbs"}}'

# List sessions
curl -sS http://localhost:8080/mbs/sessions

# Activate
curl -sS -X POST "http://localhost:8080/mbs/sessions/${ID}?action=activate"

# Update policies
curl -sS -X PATCH http://localhost:8080/mbs/sessions/${ID} \
  -H 'content-type: application/json' \
  -d '{"policies":{"qos":"mbs-gbr"}}'

# Deactivate and delete
curl -sS -X POST "http://localhost:8080/mbs/sessions/${ID}?action=deactivate"
curl -sS -X DELETE http://localhost:8080/mbs/sessions/${ID}
```

Notes

- PFCP signaling here is a placeholder: only a minimal header and UDP send are implemented. For a production implementation, build full IE encoding/decoding per 29.244 and handle responses, sequence numbers, recovery, and associations.
- The in-memory repository is not persistent. Replace with durable storage if needed.
