FROM golang:1.20.14-bullseye as build

WORKDIR /app

RUN apt-get update && apt-get install -y libolm-dev

COPY . .

# cgo must be enabled because of olm package.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-X t3.org/t3/internal/app.Version=`cat .version`" -o built/app cmd/t3/main.go

FROM golang:1.20.14-bullseye

RUN apt-get update && apt-get install -y libolm-dev ca-certificates

WORKDIR /app

COPY --from=build /app/built ./built
COPY --from=build /app/res ./res

EXPOSE 8080

ENTRYPOINT ["./built/app"]
