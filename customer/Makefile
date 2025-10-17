# Path to the environment file
ENV_FILE := .env
# Base command for running tests
GO_TEST := go test -v

# Path to the proto file
PROTO_PATH := internal/presentation/grpc/service.proto

# Base command for running protoc
PROTOC := protoc
GEN_DIR := .

.PHONY: test unit integration e2e

test: unit integration

unit:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -short ./...'

integration:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=integration ./...'

e2e:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=e2e ./...'

proto:
	$(PROTOC) --go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)
