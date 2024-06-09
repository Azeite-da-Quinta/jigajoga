# Jigajoga

In this immersive deduction board game, players travel to different zones each turn, uncovering clues to guess the identity of the other players, blending the strategic intrigue and bluff.

For your culture, jigajoga is pronounced /ʒiɡɐˈʒɔɡɐ/

## Usage

### Build for local testing

> These commands have to be run from the root of the project

```
./scripts/build.game-srv.sh
```
it will output an executable on ./game-srv/bin

### Docker compose

```
docker compose up --build -d
```

Alternative compose files are available in ./infra
