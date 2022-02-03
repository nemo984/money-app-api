postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

sqlc:
	docker run --rm -v //d/Programming/go-workspace/src/github.com/nemo984/money-app-api:/src -w //src kjconroy/sqlc generate 
