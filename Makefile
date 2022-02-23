up:
	docker-compose up -d

down:
	docker-compose down

db:
	docker exec -it money-app-api_db_1 psql -U postgres 

sqlc:
	docker run --rm -v //d/Programming/go-workspace/src/github.com/nemo984/money-app-api:/src -w //src kjconroy/sqlc generate

mockgen:
	mockgen -destination db/mock/querier.go github.com/nemo984/money-app-api/db/sqlc Querier
	mockgen -destination service/mock/service.go github.com/nemo984/money-app-api/service Service

test:
	go test -v -cover ./...

c:
	golangci-lint run

swag:
	swagger generate spec -o ./docs/swagger.yaml --scan-models

doc:
	swagger generate spec -o ./docs/swagger.yaml --scan-models && go run .
