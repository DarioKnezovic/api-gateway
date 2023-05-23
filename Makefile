# Docker image name
IMAGE_NAME := api-gateway

.PHONY: build
build:
	docker build -t $(IMAGE_NAME) .

.PHONY: run
run:
	docker run -it --rm --name $(IMAGE_NAME) -p 8080:8080 $(IMAGE_NAME)

.PHONY: hot-reload
hot-reload:
	docker run -it --rm --name $(IMAGE_NAME) -v $(PWD):/app -w /app -p 8080:8080 $(IMAGE_NAME) air

.PHONY: dev
dev:
	make build
	make hot-reload
