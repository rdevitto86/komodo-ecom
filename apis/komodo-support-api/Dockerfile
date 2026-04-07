FROM golang:1.26 AS build


WORKDIR /app

COPY komodo-support-api/go.mod komodo-support-api/go.sum ./
RUN go mod download

COPY komodo-support-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo .

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
COPY --from=build /app/internal/config/validation_rules.yaml /app/config/validation_rules.yaml
EXPOSE 7101
ENTRYPOINT ["/komodo"]
