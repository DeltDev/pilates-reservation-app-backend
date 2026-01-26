include .env
export

postgresinit:
	docker run --name postgres18 \
	-p 5433:5432 \
	-e POSTGRES_USER=root \
	-e POSTGRES_PASSWORD=password \
	-d postgres:18-alpine

postgres:
	docker exec -it postgres18 psql -U root

createdb:
	docker exec -it postgres18 createdb \
	--username=root \
	--owner=root \
	pilates

dropdb:
	docker exec -it postgres18 dropdb pilates

migrateup:
	migrate -path internal/migrations \
	-database "$(DB_SOURCE)" \
	-verbose up

migratedown:
	migrate -path internal/migrations \
	-database "$(DB_SOURCE)" \
	-verbose down
