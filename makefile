BINARY_NAME=main
BACKEND_DIR=./Backend

.PHONY: infra-up infra-down run swag build-local clean fresh

infra-up:
	docker-compose up -d db redis

infra-down:
	docker-compose stop db redis

run:
	cd $(BACKEND_DIR) && go run cmd/api/main.go

swag:
	cd $(BACKEND_DIR) && ~/go/bin/swag init -g cmd/api/main.go

build-local:
	cd $(BACKEND_DIR) && go build -o ../$(BINARY_NAME) cmd/api/main.go

clean:
	rm -f $(BINARY_NAME)
	cd $(BACKEND_DIR) && go clean

fresh:
	docker-compose down -v
	$(MAKE) infra-up