# Makefile

# Run unit tests
test:
	go test -v ./internal/command/... \
 		./internal/select/... \
 		./pkg/utils/...

# Run build binary app
build:
	go build -ldflags="-s -w" -trimpath -o dumper ./cmd/main.go

# Run the app in dev mode
run:
	go run ./cmd/main.go