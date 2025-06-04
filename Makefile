PROTO_DIR=api
PROTO_FILE=$(PROTO_DIR)/chat_v1/chat.proto
GEN_OUT_DIR=pkg	
SERVER_BIN=bin/server
CLIENT_BIN=bin/client

.PHONY: check-tools generate build-server build-client build run-server run-client clean install-lint lint run-local down

# проверяем наличие protoc и плагинов для Go
check-tools:
	@command -v protoc >/dev/null || (echo "Install protoc"; exit 1)
	@command -v protoc-gen-go >/dev/null || (echo "Install protoc-gen-go"; exit 1)
	@command -v protoc-gen-go-grpc >/dev/null || (echo "Install protoc-gen-go-grpc"; exit 1)

# генерируем Go код из proto в конкретно указанный нами место, внутри pkg
generate: check-tools
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=paths=source_relative:$(GEN_OUT_DIR) \
		--go-grpc_out=paths=source_relative:$(GEN_OUT_DIR) \
		$(PROTO_FILE)

build-server:
	go build -o $(SERVER_BIN) ./cmd/server

build-client:
	go build -o $(CLIENT_BIN) ./cmd/client

build: build-server build-client

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run -v

run-local:
	docker compose -f ./docker-compose-local.yml up --build -d

dowm: 
	docker compose -f ./docker-compose-local.yml down

# чистим сгенерированные .pb.go файлы и бинарники
clean:
	rm -rf $(GEN_OUT_DIR)/*.pb.go bin/