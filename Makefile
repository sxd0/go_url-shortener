include .env
export

# DOCKER
up:
	docker compose up --build

down:
	docker compose down -v

psql:
	docker compose exec postgres psql -U postgres link

bash:
	docker compose exec app sh


# MIGRATE
MIGRATE = migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5434/$(DB_NAME)?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-force-drop:
	$(MIGRATE) drop -f

migrate-status:
	$(MIGRATE) version


# FORMAT
format:
	@echo "Formatting all Go files"
	find . -type f -name '*.go' -exec goimports -w {} +
	go fmt ./...

build:
	@echo "Building all Go packages..."
	go build ./...

check: format build
	@echo "Format and Build Passed!!!"
