ARG DISTROLESS_TAG=debug

FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-events-api/go.mod komodo-events-api/go.sum ./
RUN go mod download

COPY komodo-events-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bootstrap ./cmd/cdc

# Lambda provided.al2023 runtime expects the handler binary to be named "bootstrap"
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /bootstrap /var/task/bootstrap
CMD ["bootstrap"]
