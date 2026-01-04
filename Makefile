BACKEND_DIR := backend
FRONTEND_DIR := frontend
GO_SOURCES := $(shell find $(BACKEND_DIR) -name '*.go' -not -path '*/vendor/*')

DOCKER_IMAGE := coderelay-runner
DOCKER_PATH := $(BACKEND_DIR)/runner

.PHONY: backend-run backend-test backend-fmt backend-lint frontend-install frontend-dev frontend-build frontend-lint compose-up build-runner-image

# Build the Docker runner image if it doesn't exist
build-runner-image:
	@if ! docker image inspect $(DOCKER_IMAGE) >/dev/null 2>&1; then \
		echo "Docker image '$(DOCKER_IMAGE)' not found, building..."; \
		docker build -t $(DOCKER_IMAGE) $(DOCKER_PATH); \
	else \
		echo "Docker image '$(DOCKER_IMAGE)' already exists"; \
	fi

backend-run: build-runner-image
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
