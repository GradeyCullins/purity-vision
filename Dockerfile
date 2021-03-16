FROM golang:latest

WORKDIR /go/src/purity-vision
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["purity-vision-filter", "-port", "8080"]
