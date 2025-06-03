install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run -v

run-local:
	docker compose -f ./docker-compose-local.yml up --build -d

dowm: 
	docker compose -f ./docker-compose-local.yml down
	
run:
	@go run cmd/main.go