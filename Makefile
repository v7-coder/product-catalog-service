.PHONY: up down restart build logs up-prod down-prod logs-prod

# === Dev ===
up:
	docker compose -f docker-compose.dev.yml up -d

down:
	docker compose -f docker-compose.dev.yml down

restart:
	docker compose -f docker-compose.dev.yml restart

build:
	docker compose -f docker-compose.dev.yml up -d --build

logs:
	docker compose -f docker-compose.dev.yml logs -f

# === Prod ===
up-prod:
	docker compose up -d

down-prod:
	docker compose down

logs-prod:
	docker compose logs -f