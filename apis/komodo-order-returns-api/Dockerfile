ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build


WORKDIR /app

COPY komodo-returns-api/go.mod komodo-returns-api/go.sum ./
RUN go mod download

COPY komodo-returns-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/server

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7062
ENTRYPOINT ["/komodo"]
