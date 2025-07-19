include .env
export

up:
	docker compose up --build

down:
	docker compose down -v

psql:
	docker compose exec postgres psql -U postgres link

bash:
	docker compose exec app sh


MIGRATE = migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5434/$(DB_NAME)?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-force-drop:
	$(MIGRATE) drop -f

migrate-status:
	$(MIGRATE) version