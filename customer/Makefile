# Path to the environment file
ENV_FILE := .env
# Base command for running tests
GO_TEST := go test -v

.PHONY: test unit integration

test: unit integration

unit:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -short ./...'

integration:
	@bash -c 'set -a; source $(ENV_FILE); set +a; $(GO_TEST) -tags=integration ./...'
