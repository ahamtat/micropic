# Build applications from sources
FROM golang:alpine AS builder
LABEL maintainer="ahamtat108@gmail.com"

# Update packages and install dependecies
RUN apk update && apk upgrade \
    && apk add ca-certificates && update-ca-certificates \
    && apk add openssl \
    && apk add git \
    && rm -rf /var/cache/apk/*

# Get sources from repository
RUN go get -v github.com/ahamtat/micropic/...
WORKDIR $GOPATH/src/github.com/ahamtat/micropic
RUN go mod vendor

# Build cache application
WORKDIR $GOPATH/src/github.com/ahamtat/micropic/cmd/cache
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/cache

# Build previewer application
WORKDIR $GOPATH/src/github.com/ahamtat/micropic/cmd/previewer
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/previewer

# Build proxy application
WORKDIR $GOPATH/src/github.com/ahamtat/micropic/cmd/proxy
RUN CGO_ENABLED=0 GOOS=linux go build -o ../../bin/proxy
