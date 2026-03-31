# BUILD_TARGET selects which cmd entrypoint to compile.
# Values: "public" (default, port 7061) | "internal" (port — shared network namespace)
ARG BUILD_TARGET=public

# debug tag includes busybox (wget) for local healthchecks.
# Override for production builds: --build-arg DISTROLESS_TAG=latest
ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

ARG BUILD_TARGET

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-order-api/go.mod komodo-order-api/go.sum ./
RUN go mod download

COPY komodo-order-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/${BUILD_TARGET}

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7061
ENTRYPOINT ["/komodo"]
