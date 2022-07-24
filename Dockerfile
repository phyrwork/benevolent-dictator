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

COPY server.go .
COPY pkg/. pkg/
RUN go build server.go


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/server .
COPY web/build dist/
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
