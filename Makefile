DB_URL := postgres://postgres:postgres@db:5432/go-api?sslmode=disable
NETWORK := heroes-go_internal

SERVICE ?=
VERSION ?=
NAME ?=

.PHONY: service-start service-stop up down restart logs ps \
        migrate-up migrate-down migrate-force migrate-version migrate-create \
        db-reset admin-keys

run:
	air

# --------------------
# Admin keys
# --------------------
admin-keys:
	mkdir -p .keys
	openssl genpkey -algorithm ed25519 -out .keys/admin_private.pem
	openssl pkey -in .keys/admin_private.pem -pubout -out .keys/admin_public.pem
# --------------------
# Docker service
# --------------------
service-start:
	sudo service docker start

service-stop:
	sudo service docker stop

# --------------------
# Containers
# --------------------
up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f $(SERVICE)

ps:
	docker compose ps

# --------------------
# Migrations
# --------------------
migrate-up:
	docker run --rm --network $(NETWORK) \
		-v $(PWD)/internal/infra/db/migrations:/migrations \
		migrate/migrate -path=/migrations -database "$(DB_URL)" up

migrate-down:
	docker run --rm --network $(NETWORK) \
		-v $(PWD)/internal/infra/db/migrations:/migrations \
		migrate/migrate -path=/migrations -database "$(DB_URL)" down 1

migrate-force:
	docker run --rm --network $(NETWORK) \
		-v $(PWD)/internal/infra/db/migrations:/migrations \
		migrate/migrate -path=/migrations -database "$(DB_URL)" force $(VERSION)

migrate-version:
	docker run --rm --network $(NETWORK) \
		-v $(PWD)/internal/infra/db/migrations:/migrations \
		migrate/migrate -path=/migrations -database "$(DB_URL)" version

migrate-create:
	migrate create -ext sql -dir internal/infra/db/migrations -seq $(NAME)

migrate-drop:
	docker run --rm -i --network $(NETWORK) \
		-v $(PWD)/internal/infra/db/migrations:/migrations \
		migrate/migrate -path=/migrations -database "$(DB_URL)" drop

# --------------------
# Dev reset
# --------------------
db-reset:
	docker compose down -v
	docker compose up -d

show-optimize:
	fieldalignment ./...

fix-optimize:
	fieldalignment -fix ./...

vulnerabilities:
	govulncheck ./...
