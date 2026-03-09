FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-communications-api/go.mod komodo-communications-api/go.sum ./
RUN go mod download

COPY komodo-communications-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./cmd/komodo-communications-api

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
EXPOSE 7081
ENTRYPOINT ["/komodo"]
