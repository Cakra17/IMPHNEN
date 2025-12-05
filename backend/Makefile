
.PHONY: run build clean migrate-up migrate-down docker-build docker-up docker-down swagger deps lint fmt

run:
	go run cmd/main.go


build:
	go build -o bin/app cmd/main.go

clean:
	rm -rf bin/

migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/migrate.sh down

docker-build:
	docker build -t imphnen-backend .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

swagger:
	swag init -g cmd/main.go -o docs/

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...
