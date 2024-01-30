#!/bin/sh
# Builder

FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/user-management-api

# Final docker image

FROM alpine:3.19

ENV SPEC_FILE_PATH=config

WORKDIR /
COPY --from=builder /app .

CMD ["/user-management-api"]