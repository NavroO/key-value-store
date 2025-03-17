BINARY=key-value-store

MAIN=cmd/server/main.go

run:
	go run $(MAIN) || true

test:
	go test ./pkg/... -v

benchmark:
	go test -bench=. -benchtime=10s ./pkg/api/
