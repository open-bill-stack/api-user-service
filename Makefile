# Завантажуємо змінні з .env у середовище
include .env
export $(shell sed 's/=.*//' .env)
# Makefile для компіляції .proto файлів за допомогою protoc

# Змінні для налаштувань
PROTO_PATH=--proto_path=proto
GO_OUT=--go_out=internal/service/grpc/proto
GO_GRPC_OUT=--go-grpc_out=internal/service/grpc/proto
GO_OPT=--go_opt=paths=source_relative
GO_GRPC_OPT=--go-grpc_opt=paths=source_relative

# Ціль за замовчуванням
all: generate

# Ціль для генерації Go файлів з .proto
generate:
	protoc $(PROTO_PATH) $(GO_OUT) $(GO_OPT) $(GO_GRPC_OUT) $(GO_GRPC_OPT) proto/user/v1/*.proto

migrate-up:
	@echo "Running migrations..."
	migrate -path ./migrations -database "$(DATABASE_URL)?sslmode=disable" up
	@echo "Migrations completed."

migrate-down:
	@echo "Running migrations..."
	migrate -path ./migrations -database "$(DATABASE_URL)?sslmode=disable" down
	@echo "Migrations completed."