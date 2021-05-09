FROM golang:alpine AS builder

WORKDIR /build

COPY . .
RUN go build -o main cmd/api/main.go

WORKDIR /dist

RUN cp /build/main .

FROM alpine

WORKDIR /app

COPY --from=builder /dist/main .
COPY db/migrations db/migrations

ENTRYPOINT ["/app/main"]