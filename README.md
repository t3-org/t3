Shield API Implemented using Golang.

__Prerequisites__

- Redis
- MongoDB

### How to run the project

- dependencies: monogDB, redis, minio, ldap(optional)
- generate a config file: `cp config.example.yaml config.yaml`
- Update config values.
- Set proper minio configs:

```bash
 # Assume the alias name is "shield-local"
mc mb shield-local/shield/p # for public files
mc policy set download shield-local/shield/p
  ```

- Run the server:

```bash
  go run ./cmd/shield/main.go server listen
```

### How to update docs?

- Install all node dependencies in the `./docs/general` directory.

```bash
yarn --cwd ./docs/general install
```

- Run docs server and then update docs: `make docs-server`

### How to deploy docs using Docker?

- build docs: `make build-docs`.
- build the image: `docker build -t shield-docs -f ./docs/general/Dockerfile ./docs/general`
- Use the image in your server. To run it locally run the following command:

```bash
docker run --rm -p 8080:80 shield-docs
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
	_ "space.org/space/internal/doc"
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
- To generate public/private keys run following commands:

```bash
ssh-keygen -f barkat_key -t rsa -b 4096 -m pem
# To get pem format for your public key:
ssh-keygen -f barkat_key.pub -e -m pem > barkat_key.pub.pem 
```

- generated public/private keys in the `secret.yaml` file in the `deloy/dev/secret.yaml` path is just a sample, genrate
  new key pairs for your production please.

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

### How to enable ldap in local env:

- run `docker-compose up -d`

- LDAP seed file is: `testdata/ldap_bootstrap.ldif` which containers multiple users.

- All LDAP users password are: `12345`

- When you import records, LDAP will generates new `entryUUID` field for them. It's the field that we use, to fetch LDAP
  records, we keep it in our database, so when you import LDAP users next time, you need to set
  `entryUUID` field for the user that you want to use to login for testing purpose (change this value manually on DB in
  your dev environment)

- run Shield.

### Deployment notes

- To run the app's server, run `{build_app} server listen` command.
- At first time run , the `shield_service_id`, `oauth_client.client_id` and `oauth_client.client_secret` should be
  empty.
- Put the app's config file in `/etc/shield/config.yaml` path.
- Run the server (deploy the app)
- To get the shield service id, oauth client id and the secret run following command in the app's container:

```bash
{built_app} oauth client info shield
```

- Put the `shield_service_id`, `oauth_client.client_id` and `oauth_client.client_secret` in the config file and restart
  the app.
- liveness, readiness endpoints:

```http request
// liveness
http://{probe_server_address}/live

// Readiness
http://{probe_server_address}/ready
```

- To just run a simple LDAP server without docker-compose file, run following command:

```bash
docker run \
        --volume "$(pwd)/testdata/ldap_bootstrap.ldif:/container/service/slapd/assets/config/bootstrap/ldif/custom/bootstrap.ldif" \
        osixia/openldap:stable --copy-service
```

- To generate `claims` param in `authorization_code`, at the `url` which redirect the user to the Shield:

```bash
echo -n '{"id_token":{"name":null,"given_name":null}}' | jq -sRr @uri
```

