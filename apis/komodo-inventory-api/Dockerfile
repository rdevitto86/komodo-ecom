ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-inventory-api/go.mod komodo-inventory-api/go.sum ./
RUN go mod download

COPY komodo-inventory-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/server

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
EXPOSE 7044
ENTRYPOINT ["/komodo"]
