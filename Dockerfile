FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /build
COPY ./main.go .
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
RUN go build -o main

FROM golang:1.22.1-alpine3.19
WORKDIR /app
COPY --from=builder /build/main .
CMD ["./main"]