# ---------- Variáveis ----------
COMPOSE_PROD = docker compose -f docker-compose.yml

# ---------- Pipeline completo ----------
.PHONY: all
all: web-build docker-build docker-up

# ---------- Prod ----------
.PHONY: docker-up docker-down docker-logs docker-build

docker-up:
	$(COMPOSE_PROD) up -d

docker-down:
	$(COMPOSE_PROD) down --volumes

docker-logs:
	$(COMPOSE_PROD) logs -f

docker-build:
	$(COMPOSE_PROD) up --build -d

# ---------- Backend Local ----------
.PHONY: api-run api-test api-build

api-run:
	go run cmd/api/main.go

api-test:
	go test ./... -v

api-build:
	go build -o bin/api ./cmd/api

# ---------- Frontend Local ----------
.PHONY: web-dev web-build

web-dev:
	cd web && npm install && npm run dev

web-build:
	cd web && npm install && npm run build