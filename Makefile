build-compose:
	docker build -t kadlab .
	docker-compose up --build

compose-up:
	docker-compose up --build

compose-down:
	docker-compose down

clirun:
	docker exec -it $(id) /bin/sh -c "/cli"