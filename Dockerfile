# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependency downloads separately from source code
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o bin/api \
    ./cmd/api/main.go

# ── Final stage ───────────────────────────────────────────────────────────────
# distroless/static: no shell, no package manager → minimal attack surface
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./api"]
