db:
	docker exec -it money-app-api_db_1

sqlc:
	docker run --rm -v //d/Programming/go-workspace/src/github.com/nemo984/money-app-api:/src -w //src kjconroy/sqlc generate

test:
	go test -cover ./..


