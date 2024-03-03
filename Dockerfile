FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY  . .
RUN go build -o coding_service cmd/main.go

FROM ubuntu:23.04 as run_stage
WORKDIR /out
COPY --from=builder /app/coding_service ./binary
CMD ["./binary"]