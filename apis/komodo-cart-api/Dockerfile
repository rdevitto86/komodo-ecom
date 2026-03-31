ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-cart-api/go.mod komodo-cart-api/go.sum ./
RUN go mod download

COPY komodo-cart-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo .

FROM gcr.io/distroless/base-debian12:${DISTROLESS_TAG}
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7043
ENTRYPOINT ["/komodo"]
