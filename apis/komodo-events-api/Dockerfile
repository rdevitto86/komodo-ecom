# BUILD_TARGET selects which cmd entrypoint to compile.
# Values: "private" (default, port 7002) | "cdc" (Lambda CDC handler)
ARG BUILD_TARGET=private
ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

ARG BUILD_TARGET

WORKDIR /app

COPY komodo-events-api/go.mod komodo-events-api/go.sum ./
RUN go mod download

COPY komodo-events-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/${BUILD_TARGET}

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7002
ENTRYPOINT ["/komodo"]
