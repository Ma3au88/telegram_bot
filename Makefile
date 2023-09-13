.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t test_telegram_bot .

build-postgres:
	docker run --name=telegram_bot -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres

start-container:
	docker run --env-file .env -p 80:80 test_telegram_bot

migrations_up:
	 migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/telegram_bot?sslmode=disable' up

migrations_down:
	 migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/telegram_bot?sslmode=disable' down