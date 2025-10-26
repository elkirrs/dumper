test:
	go test -v ./internal/command/... \
 		./internal/select/...

build:
	go build -o dumper ./cmd/main.go

run:
	go run ./cmd/main.go