run :
	go run ./cmd/main.go
postgres:
	docker-compose -f infra/docker-compose.yml up -d