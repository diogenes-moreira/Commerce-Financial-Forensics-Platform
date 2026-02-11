.PHONY: dev build test lint clean docker-up docker-down frontend-dev frontend-build

# Backend
dev:
	$(MAKE) -C backend dev

build:
	$(MAKE) -C backend build

test:
	$(MAKE) -C backend test

lint:
	$(MAKE) -C backend lint

# Frontend
frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# All
clean:
	$(MAKE) -C backend clean
	rm -rf frontend/dist frontend/node_modules
