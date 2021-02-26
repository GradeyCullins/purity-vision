run: build database
	./purity-vision-filter

build:
	go build

database:
	./start-db.sh

test:
	go test ./...

clean:
	rm purity-vision-filter
	docker stop purity-pg

.PHONY: build clean
