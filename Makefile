.PHONY: swag
swag:
	swag init -g internal/handler/http/api.go -o api/gen/openapi
