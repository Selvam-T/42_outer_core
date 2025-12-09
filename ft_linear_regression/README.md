## 1. Docker Setup

### Go Development Container (go-service)

- Mounts host folder ./workspace to /workspace for live coding.
- Provides Go environment and tools (go, git, nano, etc.).

### Docker Compose

- docker compose up -d starts container.
- docker compose exec go-service bash enters Go environment.

### Makefile Commands

- make build → build images
- make up → start containers
- make go-shell → enter Go container
- make logs → view logs for all services
- make down → stop containers

- make launch → launch go container, bash into container
- make exit  → down, rm images

