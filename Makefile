build-compose:
	docker build -t kadlab .
	docker-compose up -d --build

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down