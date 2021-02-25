db_running=$(docker ps)

build:
	go build

database:
	./start-db.sh

run: build database
	./purity-vision-filter

test:
	go test ./...

.PHONY: build
