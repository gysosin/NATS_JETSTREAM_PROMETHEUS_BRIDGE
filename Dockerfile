FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o exporter ./cmd/exporter

FROM alpine:latest
COPY --from=builder /app/exporter /usr/local/bin/exporter
COPY config.json /etc/exporter/config.json
EXPOSE 2112
ENTRYPOINT ["/usr/local/bin/exporter"]
