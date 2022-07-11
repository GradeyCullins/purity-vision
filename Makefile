SOURCES := $(shell find . -type f -name '*.go')
TARGET = purity-vision
.DEFAULT_GOAL: $(TARGET)
TAG = latest

.PHONY: docker-run run test docker-stop clean

docker-run: $(TARGET)
	docker-compose up --detach

run: $(TARGET)
	./start-db.sh
	./${TARGET}

$(TARGET): $(SOURCES) Dockerfile .envrc
	GOOS=linux GOARCH=amd64 go build -o ${TARGET}
	docker build -t ${TARGET}:${TAG} .

test:
	PURITY_DB_HOST="localhost" go test ./...

docker-stop:
	docker-compose down

clean:
	rm ${NAME}
	docker stop purity-pg
