ARG GO_VERSION=1.22.3

# build stage
FROM golang:${GO_VERSION} AS build-stg

WORKDIR /build
COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY ./game-srv ./game-srv
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o cabrito -ldflags '-w -extldflags "-static"' ./game-srv/main.go

# test stage
FROM build-stg AS test-stg

RUN go test -v ./...

# run stage
FROM alpine

WORKDIR /
COPY --from=build-stg /build/cabrito /cabrito

EXPOSE 80

# TODO certificates and stuff
ENTRYPOINT [ "/cabrito", "serve" ]
