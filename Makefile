build-compose:
	docker build -t kadlab .
	docker-compose up --build

compose-up:
	docker-compose up --build

compose-down:
	docker-compose down

run:
	docker exec -it $(id) /bin/sh -c "go run cli/cli.go"