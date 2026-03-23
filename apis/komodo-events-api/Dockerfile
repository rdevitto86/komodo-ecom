ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-events-api/go.mod komodo-events-api/go.sum ./
RUN go mod download

COPY komodo-events-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/server

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
EXPOSE 7002
ENTRYPOINT ["/komodo"]
