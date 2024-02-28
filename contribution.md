### How to generate api docs?

- Extract new routes:

```bash
make extract-api-docs
```

- Update generated document for new routes in the `internal/router/api/doc/api_docs.go` file.

- Generate openapi docs

```bash
 make build-api-docs
```

- Run the api docs server: `make serve-api-docs`

- For production run docs on a server.

__Notes for developers__

- To sync your packages with the `vendor` dir run `go mod vendor` command.
