
FROM golang:1.21.3-alpine3.18 as builder

WORKDIR /opt

COPY  ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go build -o /opt/build/server /opt/cmd/server

FROM alpine:3.18.4

COPY --from=builder /opt/build/server /usr/local/bin/server

CMD server