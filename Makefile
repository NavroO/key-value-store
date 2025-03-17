BINARY=key-value-store

MAIN=cmd/server/main.go

run:
	go run $(MAIN)

test:
	go test ./pkg/... -v
