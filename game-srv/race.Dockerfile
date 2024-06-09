# TODO not sure it works 100%
ARG GO_VERSION=1.22.3

FROM golang:${GO_VERSION}

WORKDIR /build
COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./game-srv ./game-srv
RUN GORACE="log_path=./report" CGO_ENABLED=1 GOOS=linux \
    go build --race -o cabrito -ldflags '-w -extldflags "-static"' ./game-srv/main.go

EXPOSE 80

ENTRYPOINT [ "./cabrito", "serve" ]
