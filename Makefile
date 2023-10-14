MIGRATE_DSN=host=localhost port=6432 user=postgres password=password database=be

install-deps:
	GOBIN=$(shell pwd)/bin go install github.com/pressly/goose/v3/cmd/goose@v3.11.2

run-db:
	cd .docker && sudo docker-compose -f docker-compose.local.yml up -d postgres && cd ..
	sleep 5
	$(shell pwd)/bin/goose -dir migrations/ -allow-missing postgres "$(MIGRATE_DSN)" up

docker-down:
	cd .docker;\
	sudo docker-compose -f docker-compose.local.yml down;\
	cd ..

run:
	go run cmd/main.go -address localhost:3001

test:
