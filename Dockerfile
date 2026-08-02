FROM golang:1.23-alpine3.21 AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/iknow-tools .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates coreutils

WORKDIR /app
COPY --from=builder /out/iknow-tools ./iknow-tools
COPY cron.sh ./cron.sh
COPY docker/collector-entrypoint.sh ./docker/

RUN chmod +x ./iknow-tools ./cron.sh ./docker/*.sh

ENV IKNOW_DATA_DIR=/data

ENTRYPOINT ["/app/docker/collector-entrypoint.sh"]
