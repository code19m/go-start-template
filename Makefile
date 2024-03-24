# To use all commands in this Makefile, you need to install the following tools:
# - golangci-lint
# - docker
# - goose
# - swag
# - air


.PHONY: run
run:
	air app --http_addr=localhost:8081


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: swag
swag:
	swag init -g internal/handler/http/api.go -o api/gen/openapi


.PHONY: devup
devup:
	docker-compose -f docker-compose.local.yaml up -d --build


.PHONY: devdown
devdown:
	docker-compose -f docker-compose.local.yaml down


.PHONY: devexec
devexec:
	docker-compose -f docker-compose.local.yaml exec -it $(name) bash


.PHONY: migrate-add
migrate-add:
	goose -dir=migrations create $(name) sql


.PHONY: migrate-up
migrate-up:
	# Set GOOSE_DRIVER GOOSE_DBSTRING variables before using
	goose -dir=migrations up


.PHONY: migrate-down
migrate-down:
	# Set GOOSE_DRIVER GOOSE_DBSTRING variables before using
	goose -dir=migrations down
