# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder

MAINTAINER "2411167570@qq.com"

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace


COPY .. .

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o main ./main.go

ENV TZ=Asia/Shanghai
ENTRYPOINT ["/main"]

