FROM golang:alpine AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .

RUN go build -o main cmd/api/main.go

WORKDIR /dist

RUN cp /build/main .

FROM alpine

WORKDIR /app

COPY --from=builder /dist/main .
COPY db/migrations db/migrations

ENTRYPOINT ["/app/main"]