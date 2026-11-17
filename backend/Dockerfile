FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /seed ./cmd/seed
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /api /api
COPY --from=builder /worker /worker
COPY --from=builder /seed /seed
COPY --from=builder /migrate /migrate
COPY --from=builder /app/migrations ./migrations

ENV TZ=Asia/Kolkata

EXPOSE 8080

ENTRYPOINT ["/api"]