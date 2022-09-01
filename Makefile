.PHONY: update-repos
update-repos:
	go mod tidy
	go mod vendor

.PHONY: layers
layers:
	go run ./internal/app/layer_generator/main.go tracing
	go run ./internal/store/layer_generator/main.go

.PHONY: listen
listen:
	go run ./cmd/space/main.go server listen

.PHONY: extract-api-docs
extract-api-docs:
	go run ./cmd/space/main.go openapi extract

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

