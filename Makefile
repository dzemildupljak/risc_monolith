include .env
include .env.heroku

run:
	go run server/*.go

watch:
	reflex -r 'server/.*\.go$$' -s go run server/*.go 



docker-build:
	docker build -f ./docker/Dockerfile .

docker-compose-dev-build:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.docker.dev build

docker-compose-dev-up:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.docker.dev up

docker-compose-dev-down:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.docker.dev down

docker-compose-dev-config:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.docker.dev config

docker-develop-migrateup:
	docker run -v /home/dzemil/Projects/golang/risc_monolith/server/db/postgres/migrations:/server/db/postgres/migrations --network host migrate/migrate -path=/server/db/postgres/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

docker-develop-migratedown:
	docker run -v /home/dzemil/Projects/golang/risc_monolith/server/db/postgres/migrations:/server/db/postgres/migrations --network host migrate/migrate -path=/server/db/postgres/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down -all

docker-compose-prod-build:
	docker-compose -f docker/docker-compose.yml --env-file ./.env build

docker-compose-prod-up:
	docker-compose -f docker/docker-compose.yml --env-file ./.env up

docker-compose-prod-up-build:
	docker-compose -f docker/docker-compose.yml --env-file ./.env up --build

docker-compose-prod-down:
	docker-compose -f docker/docker-compose.yml --env-file ./.env down

docker-compose-prod-config:
	docker-compose -f docker/docker-compose.yml --env-file ./.env config


heroku-logs:
	heroku logs --tail -a serene-fortress-45917

heroku-container-push: 
	cd docker/ && heroku container:push web --app serene-fortress-45917 --context-path ../ && cd ..
heroku-container-push-migrate: 
	cd docker/ && heroku container:push web --app serene-fortress-45917 --context-path ../ && cd ..
	make migrate-up-heroku
heroku-container-release:
	cd docker/ && heroku container:release web --app serene-fortress-45917 && cd ..



create-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir server/db/postgres/migrations -seq $$migration_name
migrate-up:
	migrate -path server/db/postgres/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose up

migrate-down:
	migrate -path server/db/postgres/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" -verbose down

migrate-up-heroku:
	migrate -path server/db/postgres/migrations -database "$(HEROKU_DB_URI)" -verbose up

migrate-down-heroku:
	migrate -path server/db/postgres/migrations -database "$(HEROKU_DB_URI)" -verbose down



generate-rsa-private-key-access:
	openssl genrsa -out access-private.pem 2048
generate-rsa-public-key-access:
	openssl rsa -in access-private.pem -outform PEM -pubout -out access-public.pem
generate-rsa-private-key-refresh:
	openssl genrsa -out refresh-private.pem 2048
generate-rsa-public-key-refresh:
	openssl rsa -in refresh-private.pem -outform PEM -pubout -out refresh-public.pem

generate-new-rsa-keys:
	make generate-rsa-private-key-access
	make generate-rsa-private-key-refresh
	make generate-rsa-public-key-access
	make generate-rsa-public-key-refresh

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.json --scan-models
