
LDFLAGS = -extldflags=-Wl,-ld64

.PHONY: update-repos
update-repos:
	go mod tidy
	go mod vendor

.PHONY: layers
layers:
	go run -ldflags='$(LDFLAGS)' ./internal/app/layer_generator/main.go tracing
	go run -ldflags='$(LDFLAGS)' ./internal/store/layer_generator/main.go

.PHONY: listen
listen:
	go run -ldflags='$(LDFLAGS)' ./cmd/t3/main.go server listen

.PHONY: extract-api-docs
extract-api-docs:
	go run -ldflags='$(LDFLAGS)' ./cmd/t3/main.go openapi extract

.PHONY: build-api-docs
build-api-docs:
	swagger generate spec -o ./docs/api/api_docs.json

.PHONY: api-docs-server
api-docs-server:
	 serve -l 1000 ./docs/api

.PHONY: publish-api-docs
publish-api-docs:
	vercel ./docs/api --prod

.PHONY: docs-server
docs-server:
	yarn --cwd ./docs/general docs:dev .

.PHONY: build-docs
build-docs:
	yarn --cwd ./docs/general docs:build .

.PHONY: publish-docs
publish-docs:
	vercel ./docs/general --prod

.PHONY: install-lint
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: gen
gen:
	go generate ./...
