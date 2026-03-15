-include .env
export

APP_NAME := api
BIN_DIR := bin

.PHONY: run build docker-up

run:
	go run ./cmd/api

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

docker-up:
	docker compose up --build
