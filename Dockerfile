FROM golang:latest AS build

WORKDIR /go/src/purity-vision
COPY . .

RUN go get -d -v ./...

# Add new stage to cache Go dependency download.
FROM build

RUN go install -v ./...

EXPOSE 8080

CMD ["purity-vision", "-port", "8080"]
