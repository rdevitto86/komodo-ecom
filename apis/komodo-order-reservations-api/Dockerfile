FROM golang:1.26 AS build


WORKDIR /app

COPY komodo-order-reservations-api/go.mod komodo-order-reservations-api/go.sum ./
RUN go mod download

COPY komodo-order-reservations-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
EXPOSE 7063
ENTRYPOINT ["/komodo"]
