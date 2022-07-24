# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18 as build

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY tools.go .
RUN go mod download

COPY cmd/. cmd/
COPY pkg/. pkg/
RUN cd cmd/api && go build


##
## Deploy
##
FROM ubuntu:20.04

WORKDIR /

COPY --from=build /app/cmd/api/api /usr/local/bin/api

EXPOSE 8080

ENTRYPOINT /usr/local/bin/api
