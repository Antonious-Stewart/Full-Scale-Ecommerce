.DEFAULT_GOAL := run

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run main.go