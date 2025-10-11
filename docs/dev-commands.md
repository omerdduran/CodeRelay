# Quick Dev Commands

- `make backend-run`: start the Go API (`go run ./cmd/api`).
- `make backend-test`: execute backend tests (`go test ./...`).
- `make backend-fmt`: format Go sources (`gofmt -w`).
- `make backend-lint`: run `golangci-lint` against the backend.
- `make frontend-install`: install Node dependencies.
- `make frontend-dev`: launch the Next.js dev server.
- `make frontend-build`: generate a production Next.js bundle.
- `make frontend-lint`: run ESLint checks.
- `make compose-up`: build and launch api/frontend/sqlite via Docker Compose.
