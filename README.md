T3(The Ticket Tracker) tracks the tickets (e.g., grafana alerts).


![alt text](./docs/screenshots/home.png)

__Project status__: Under development. it's not stable yet.

__Prerequisites__

- Postgres

### How to run the project

- generate config file: `cp config.example.yaml config.yaml`
- generate channels file: `cp chanenls.example.yaml channels.yaml`
- Update config and channels values.
- Run the server: `go run ./cmd/t3/main.go server listen`
- Run the UI server by following its steps [here](./webui/README.md)


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

- To run the app's server, run `{built_app} server listen` command.
- liveness, readiness endpoints:

```text
// liveness
http://{probe_server_address}/live

// Readiness
http://{probe_server_address}/ready
```

### TODO
- Read the [next steps](./docs/general/next_steps.md)
- Write tests
