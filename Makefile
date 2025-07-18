up:
	docker compose up --build

down:
	docker compose down -v

migrate:
	docker compose run --rm migrate

psql:
	docker compose exec postgres psql -U postgres link

bash:
	docker compose exec app sh
