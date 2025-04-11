# Path to the environment file
ENV_FILE := .env

# Base command for running tests
GO_TEST := go test -v

# Path to the proto file
PROTO_PATH := internal/presentation/grpc/service.proto

# Base command for running protoc
PROTOC := protoc

.PHONY: test unit integration proto

test: unit integration

unit:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -short ./...'

integration:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=integration ./...'

proto:
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)
