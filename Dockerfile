FROM golang:1.23-bookworm AS builder
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/jhldy .

FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/jhldy /app/jhldy
COPY config.json /app/config.json
COPY web /app/web
COPY templates /app/templates
COPY downloads/app-release.apk /app/downloads/app-release.apk
ENV JH_LDY_ROOT=/app
EXPOSE 5175
HEALTHCHECK --interval=15s --timeout=3s --retries=3 \
  CMD curl -fsS http://127.0.0.1:5175/health || exit 1
CMD ["/app/jhldy"]
