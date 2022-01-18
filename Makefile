run:
	go run server/*.go

watch:
	reflex -r 'server/.*\.go$$' -s go run server/*.go 

docker-build:
	docker build -f ./docker/Dockerfile .
w
docker-compose-dev-build:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.dev build

docker-compose-dev-up:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.dev up

docker-compose-dev-down:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.dev down

docker-compose-dev-config:
	docker-compose -f docker/docker-compose.dev.yml --env-file ./.env.dev config


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
	cd docker/;
	heroku container:push web --app serene-fortress-45917 --context-path ../;
	cd ..
heroku-container-release:
	cd docker/;
	heroku container:release web --app serene-fortress-45917;
	cd ..
