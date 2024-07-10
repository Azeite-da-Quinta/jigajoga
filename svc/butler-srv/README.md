# Butler server

TODO

## Usage

### Locally

```
# app welcome
go run main.go

# run the server
go run main.go serve --level=DEBUG --port=8080
```

### Docker

> These commands have to be run from the root of the project

```
# build butler-srv image
docker build -t butler-srv:latest -f svc/butler-srv/Dockerfile .

# run the image
docker run butler-srv
```