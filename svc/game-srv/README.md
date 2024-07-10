# Game server

This server program is designed to host an online version of an immersive deduction board game.

## Usage

### Locally

```
# app welcome
go run main.go

# run the server
go run main.go serve --level=DEBUG --port=8080

# run the client
go run main.go client --host=127.0.0.1:8080
```

### Docker

> These commands have to be run from the root of the project

```
# build game-srv image
docker build -t game-srv:latest -f svc/game-srv/Dockerfile .

# run the image
docker run game-srv
```