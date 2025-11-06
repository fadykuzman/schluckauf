FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
  gcc \
  libc6-dev \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
  -ldflags="-w -s" \
  -o /app/dup-reviewer ./cmd/dup-reviewer


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  sqlite3 \
  wget \
  && rm -rf /var/lib/apt/lists/*

ARG CZKAWKA_CLI_VERSION=10.0.0
ARG CZKAWKA_CLI_SHA256=b261aba0ca0b1d99d450949be22f9ae172750fe13dc9b40a32209fc8db0fc159


RUN wget --progress=dot:giga -O /usr/local/bin/czkawka_cli \
  https://github.com/qarmin/czkawka/releases/download/${CZKAWKA_CLI_VERSION}/linux_czkawka_cli_x86_64 && \
  echo "${CZKAWKA_CLI_SHA256}  /usr/local/bin/czkawka_cli" > /tmp/checksum && \
  sha256sum -c /tmp/checksum && \
  chmod +x /usr/local/bin/czkawka_cli && \
  rm /tmp/checksum

RUN groupadd -g 1000 appuser && \
  useradd -u 1000 -g appuser -m appuser

RUN mkdir -p /photos /data /trash /scans && \
  chown -R appuser:appuser /photos /data /trash /scans

COPY --from=builder /app/dup-reviewer /usr/local/bin/
COPY --chown=appuser:appuser web /app/web

USER appuser

ENV DATABASE_PATH=/data/duplicates.db \
  TRASH_DIR=/trash \
  SCANS_DIR=/scans \
  PHOTOS_DIR=/photos

WORKDIR /app

EXPOSE 8080

CMD ["dup-reviewer"]


