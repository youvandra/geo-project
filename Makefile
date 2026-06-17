.PHONY: blog blog-dev build-cli build-backend docker-up clean

# Blog
blog:
	cd blog && hugo

blog-dev:
	cd blog && hugo server -D

# CLI
build-cli:
	cd cli && mkdir -p ../bin && go build -o ../bin/geo .

# Backend
build-backend:
	cd backend && go build -o ../bin/geo-backend ./cmd/server/

run-backend:
	cd backend && go run ./cmd/server/

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Misc
clean:
	rm -rf bin/ blog/public/

all: build-cli build-backend
