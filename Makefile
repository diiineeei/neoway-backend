swagger: ## Gera documentação Swagger
	@echo "Gerando documentação..."
	@cd cmd/api && PATH=$(PATH):$(shell go env GOPATH)/bin swag init -g main.go -o ../../docs
	@echo "Swagger gerado em docs/"

build:
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v --rmi all

test:
	go test -v -race ./...

rebuild: swagger down build

