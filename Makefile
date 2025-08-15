run:
	go run .

build:
	go build .

tidy:
	go mod tidy

migrate:
	go run /database/schemas.go
