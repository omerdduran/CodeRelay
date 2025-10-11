BACKEND_DIR := backend
FRONTEND_DIR := frontend
GO_SOURCES := $(shell find $(BACKEND_DIR) -name '*.go' -not -path '*/vendor/*')

.PHONY: backend-run backend-test backend-fmt backend-lint frontend-install frontend-dev frontend-build frontend-lint compose-up

backend-run:
	@cd $(BACKEND_DIR) && go run ./cmd/api

backend-test:
	@cd $(BACKEND_DIR) && go test ./...

backend-fmt:
	@gofmt -w $(GO_SOURCES)

backend-lint:
	@cd $(BACKEND_DIR) && golangci-lint run ./...

frontend-install:
	@cd $(FRONTEND_DIR) && npm install

frontend-dev:
	@cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	@cd $(FRONTEND_DIR) && npm run build

frontend-lint:
	@cd $(FRONTEND_DIR) && npm run lint

compose-up:
	@docker compose up --build
