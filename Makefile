# Makefile

# Run unit tests
test:
	go test -v ./internal/command/... \
 		./internal/select/...

# Run build binary app
build:
	go build -o dumper ./cmd/main.go

# Run the app in dev mode
run:
	go run ./cmd/main.go