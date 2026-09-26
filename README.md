# Миграции

## Пример создания миграции

migrate create -ext sql -dir migrations -seq add_slug_to_products

## Пример накатить миграции

migrate -path ./migrations -database "postgres://postgres:postgres@mp-pgbouncer:5432/mp_pcs?sslmode=disable" up
