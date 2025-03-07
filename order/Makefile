.PHONY: test unit integration

test: unit integration

unit:
	@bash -c 'set -a; source .env; set +a; go test -v -short ./...'

integration:
	@bash -c 'set -a; source .env; set +a; go test -v -tags=integration ./...'
