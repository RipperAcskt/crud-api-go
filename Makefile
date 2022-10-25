build:
	go build -o ./bin/service ./cmd/main.go

run:
	./bin/service

go: build run
