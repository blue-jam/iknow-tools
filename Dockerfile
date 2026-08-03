FROM golang:1.23-alpine3.21 AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/iknow-tools .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates coreutils tzdata

RUN addgroup -S -g 10001 iknow \
    && adduser -S -D -H -u 10001 -G iknow iknow \
    && mkdir -p /data \
    && chown iknow:iknow /data

WORKDIR /app
COPY --from=builder /out/iknow-tools ./iknow-tools
COPY cron.sh ./cron.sh
COPY docker/collector-entrypoint.sh ./docker/

RUN chmod +x ./iknow-tools ./cron.sh ./docker/*.sh

ENV IKNOW_DATA_DIR=/data

USER iknow:iknow

ENTRYPOINT ["/app/docker/collector-entrypoint.sh"]
