# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /git-mirror-action ./cmd/main.go

# Runtime stage
FROM alpine:3.23

RUN apk add --no-cache git openssh-client

COPY --from=builder /git-mirror-action /usr/local/bin/git-mirror-action

ENTRYPOINT ["git-mirror-action"]
