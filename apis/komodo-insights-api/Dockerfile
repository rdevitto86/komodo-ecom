ARG BUILD_TARGET=public
ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build
ARG BUILD_TARGET

WORKDIR /app

COPY komodo-insights-api/go.mod komodo-insights-api/go.sum ./
RUN go mod download

COPY komodo-insights-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/${BUILD_TARGET}

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
EXPOSE 7111
ENTRYPOINT ["/komodo"]
