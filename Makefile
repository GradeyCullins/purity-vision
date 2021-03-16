SOURCES := $(shell find ./ -name '*.go')
NAME = purity-vision
TAG = latest

run: docker-build
	docker-compose up --detach

docker-build: build Dockerfile
	docker build -t ${NAME}:${TAG} .

build: $(SOURCES)
	go build

local: database

	./purity-vision-filter -port 8080

database:
	./start-db.sh

test:
	go test ./...

stop:
	docker-compose down

clean:
	rm purity-vision-filter
	docker stop purity-pg

.PHONY: clean docker-build run local database test down stop
