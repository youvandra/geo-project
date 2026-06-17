.PHONY: build-cli build-backend run-backend dev docker-up docker-down clean all

# CLI
build-cli:
	cd cli && mkdir -p ../bin && go build -o ../bin/geo .

# Backend
build-backend:
	cd backend && go build -o ../bin/geo-backend ./cmd/server/

run-backend:
	cd backend && go run ./cmd/server/

# Dev — run both with hot-reload (requires nodemon/entr)
dev:
	@echo "To run both:"
	@echo "  make run-backend  (dashboard at :8080)"
	@echo "  make build-cli    (CLI tools)"

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Misc
clean:
	rm -rf bin/

all: build-cli build-backend
