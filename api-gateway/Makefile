# Base command to generate swagger docs
SWAGGER := swag

.PHONY: swagger

swagger:
	${SWAGGER} init -g ./cmd/api-gateway/main.go -o ./internal/adapter/input/api/docs
