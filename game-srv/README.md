# Game server

This server program is designed to host an online version of an immersive deduction board game.

## Usage

### Locally

```
# app welcome
go run main.go

# run the client
go run main.go serve

# run the client
go run main.go client
```

### Docker

> These commands have to be run from the root of the project

```
# build game-srv image
docker build -t game-srv:latest -f game-srv/Dockerfile .

# run the image
docker run game-srv
```