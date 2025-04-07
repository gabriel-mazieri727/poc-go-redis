.PHONY: run build docker-build docker-up docker-down docker-logs

run:
	wgo run .

build:
	go build -o main .

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f