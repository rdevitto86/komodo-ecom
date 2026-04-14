ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build


WORKDIR /app

COPY komodo-user-wishlist-api/go.mod komodo-user-wishlist-api/go.sum ./
RUN go mod download

COPY komodo-user-wishlist-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/public

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7053
ENTRYPOINT ["/komodo"]
