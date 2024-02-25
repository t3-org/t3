T3(The Ticket Tracker) tracks the tickets (e.g., grafana alerts).

__Project status__: Under development. it's not stable yet.

__Prerequisites__

- Postgres

### How to run the project

- generate a config file: `cp config.example.yaml config.yaml`
- Update config values.
- Run the server: `go run ./cmd/t3/main.go server listen`
- Run the UI server following steps [here](./webui/README.md)

### How to update docs?

- Install all node dependencies in the `./docs/general` directory.

```bash
yarn --cwd ./docs/general install
```

- Run docs server and then update docs: `make docs-server`

### How to deploy docs using Docker?

- build docs: `make build-docs`.
- build the image: `docker build -t t3-docs -f ./docs/general/Dockerfile ./docs/general`
- Use the image in your server. To run it locally run the following command:

```bash
docker run --rm -p 8080:80 t3-docs
# and then open localhost:8080 on your browser.
```

### How to deploy docs using vercel?

- build docs: `make build-docs`.
- built docs are in the `./docs/general/build` path.
- Now you can deploy them. The next steps are instruction to deploy docs into vercel server:
- Change directory to `./docs/general` and initial a vercel project there.

```bash
cd ./docs/general
# In the initialization process if you want to
# deploy to already created vercel project, so link it.
# Otherwise set the deploy path to root directory to `dist`
# and deploy it.
vercel
```

- After initialization, next times just deploy it to production:

```bash
make publish-docs
```

### How to generate api docs?

- Extract new routes:

```bash
go run main.go openapi extract
```

- Update generated docs for new routes in the `internal/doc/api_docs.go` file.

- Add the docs golang package to the `main.go` file.

```go
package main

import (
	_ "t3.org/t3/internal/doc"
)
```

- Generate openapi docs using `swagger`

```bash
 swagger generate spec -o ./docs/api/api_docs.json
```

- Run a simple server in the `./docs/api` directory to server api docs locally (e.g., `serve -l 1000`).
- For production run docs on a server.

__Notes for developers__

- To sync your packages with the `vendor` dir run `go mod vendor` command.

### How to enable jaeger in local env

- Run jaeger:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.26
```

- Update app's Jaeger config:

```yaml
open_telemetry:
  tracing:
    noop_tracer: false # set true to disable jaeger
    jaeger_addr: "http://localhost:14268/api/traces"
    always_sample: true # in production set to false.
```

### Deployment notes

- To run the app's server, run `{build_app} server listen` command.
- liveness, readiness endpoints:

```text
// liveness
http://{probe_server_address}/live

// Readiness
http://{probe_server_address}/ready
```

### Notes about tests

- Right now we are using `defer testbox.Global().TeardownIfPanic()` per each test to make sure we'll clean up
  everything in situation that our test panic. Unfortunately we can not use `t.Cleanup()` function to do that, because
  when it panics, `t.Failed()` doens't returns `true` in the cleanup function,
  [but it's rolled forward to go1.20](https://github.com/golang/go/issues/49929), in go 1.20 update it please. please
  note `t.Failed()` returns `true` in all situations that our test is failed, not just a panic. so check that this is
  what we want or not (or it would be great if Go add `Cleanup` to the `testing.M`).

### HTTP request samples:

```bash
xh -v post :4000/api/v1/tickets fingerprint="h" is_firing:=true started_at:=1 level="low" description="a test alert" webhook:='{"channel":"matrix","channel_id":"!sGJfLhjEueOpYkVKdz:matrix.org"}'
```

### TODO:

- Write query filter functionality at server-side.
- Write tests
- QA
