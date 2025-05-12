# Завантажуємо змінні з .env у середовище
include .env
export $(shell sed 's/=.*//' .env)
# Makefile для компіляції .proto файлів за допомогою protoc

# Змінні для налаштувань
PROTOC=protoc
GO_OUT=--go_out=.
GO_GRPC_OUT=--go-grpc_out=.
GO_OPT=--go_opt=paths=source_relative
GO_GRPC_OPT=--go-grpc_opt=paths=source_relative

# Шлях до .proto файлу
PROTO_FILE=grpc/grpc.proto

# Ціль за замовчуванням
all: generate

# Ціль для генерації Go файлів з .proto
generate:
	find . -type f -name "*.proto" -exec $(PROTOC) $(GO_OUT) $(GO_OPT) $(GO_GRPC_OUT) $(GO_GRPC_OPT) {} +

migrate-up:
	@echo "Running migrations..."
	migrate -path ./migrations -database "$(DATABASE_URL)?sslmode=disable" up
	@echo "Migrations completed."

migrate-down:
	@echo "Running migrations..."
	migrate -path ./migrations -database "$(DATABASE_URL)?sslmode=disable" down
	@echo "Migrations completed."