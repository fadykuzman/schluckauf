FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
  -ldflags="-w -s" \
  -o /app/dup-reviewer ./cmd/dup-reviewer


FROM alpine:3.22

RUN apk add --no-cache \
  ca-certificates=20250911-r0 \
  sqlite=3.49.2-r1 \
  wget=1.25.0-r1

ARG CZKAWKA_CLI_VERSION=10.0.0
ARG CZKAWKA_CLI_SHA256=b261aba0ca0b1d99d450949be22f9ae172750fe13dc9b40a32209fc8db0fc159


RUN wget --progress=dot:giga \
  https://github.com/qarmin/czkawka/releases/download/${CZKAWKA_CLI_VERSION}/linux_czkawka_cli.tar.xz && \
  echo "${CZKAWKA_CLI_SHA256}  linux_czkawka_cli.tar.xz" > /tmp/checksum && \
  sha256sum -c /tmp/checksum && \
  tar -xf linux_czkawka_cli.tar.xz && \
  mv linux_czkawka_cli/czkawka_cli /usr/local/bin/ && \
  chmod +x /usr/local/bin/czkawka_cli && \
  rm -rf linux_czkawka_cli linux_czkawka_cli.tar.xz

RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser

RUN mkdir -p /photos /data /trash /scans && \
  chown -R appuser:appuser /photos /data /trash /scans

COPY --from=builder /app/dup-reviewer /usr/local/bin/
COPY --chown=appuser:appuser web /app/web

USER appuser

WORKDIR /app

EXPOSE 8080

CMD ["dup-reviewer"]


