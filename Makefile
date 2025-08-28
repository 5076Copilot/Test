APP=mbsmf
PKG=./...

.PHONY: build run test lint

build:
	GOFLAGS="-trimpath" CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/mbsmf

run:
	MB_SMF_HTTP_ADDR=:8080 go run ./cmd/mbsmf

test:
	go test -v $(PKG)

