
FROM golang:bookworm AS builder

WORKDIR /opt

COPY  ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go build -o /opt/build/client /opt/cmd/client

FROM ubuntu:22.04

COPY --from=builder /opt/build/client /usr/local/bin/client
