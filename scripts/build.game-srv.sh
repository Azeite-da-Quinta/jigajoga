cd ./game-srv
CGO_ENABLED=0 GOOS=linux \
    go build -o bin/game-srv -ldflags '-w -extldflags "-static"' main.go
