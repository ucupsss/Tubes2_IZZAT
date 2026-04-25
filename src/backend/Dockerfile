FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/backend ./src/backend

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/backend ./src/backend

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /out/backend /app/backend

EXPOSE 5175

CMD ["/app/backend"]
