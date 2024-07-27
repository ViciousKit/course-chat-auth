include local.env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATAION_DIR:=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN:="host=localhost user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DATABASE_NAME) port=$(PG_PORT)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-server-api

generate-server-api:
	mkdir -p generated
	protoc -I proto proto/auth_v1/auth.proto \
	--go_out=./generated \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=./generated \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \

bin-build:
	GOOS=linux GOARCH=amd64 go build -o auth_service_linux cmd/main.go

docker-build:
	docker buildx build --no-cache --platform linux/amd64 -t auth_service:v0.0.1

local-migration-status:
	bin/goose -dir $(LOCAL_MIGRATAION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	bin/goose -dir $(LOCAL_MIGRATAION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	bin/goose -dir $(LOCAL_MIGRATAION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

compose-up-dev:
	export ENV_FILE=local.env
	docker compose --profile dev --env-file local.env up -d --build

compose-down-dev:
	export ENV_FILE=local.env
	docker compose --profile dev --env-file local.env down