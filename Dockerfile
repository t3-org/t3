FROM golang:1.18.4-alpine as build

WORKDIR /app

RUN apk add --no-cache olm

COPY . .

# cgo must be enabled because of olm package.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-X t3.org/t3/internal/app.Version=`cat .version`" -o built/app cmd/t3/main.go

FROM golang:1.18.4-alpine

RUN apk add --no-cache olm ca-certificates

WORKDIR /app

COPY --from=build /app/built ./built
COPY --from=build /app/res ./res

EXPOSE 80

ENTRYPOINT ["./built/app"]
