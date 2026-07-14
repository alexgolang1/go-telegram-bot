APP_NAME=tg-bot
BUILD_DIR=bin

.PHONY: 

run:
	go run main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

clean:
	rm -rf $(BUILD_DIR)

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose down
	docker compose up -d --build