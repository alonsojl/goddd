FROM golang:1.19-alpine

WORKDIR /home/alonsojl
COPY Makefile .
RUN apk update && apk add --no-cache build-base && apk add curl
RUN make tools

