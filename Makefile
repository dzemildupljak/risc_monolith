run:
	go run server/*.go

watch:
	reflex -r 'server/.*\.go$$' -s go run server/*.go 

docker-build:
	docker build -f ./docker/Dockerfile .