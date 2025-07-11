FROM golang:1.24-alpine AS build

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o main cmd/server.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates curl sqlite tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

RUN mkdir -p /app/data /app/assets && \
    chown -R appuser:appgroup /app

# Copy files without changing ownership initially
COPY --from=build /app/main /app/main
COPY assets /app/assets
COPY scripts/healthcheck.sh /app/healthcheck.sh

# Set minimal required permissions and ownership
RUN chmod 755 /app/main && \
    chmod 755 /app/healthcheck.sh && \
    chown appuser:appgroup /app/main /app/healthcheck.sh && \
    chown -R appuser:appgroup /app/assets && \
    chmod -R 644 /app/assets

USER appuser


EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD ["/app/healthcheck.sh"]

CMD ["./main"]
