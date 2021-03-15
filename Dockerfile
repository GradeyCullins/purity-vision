FROM golang:latest

# Install direnv for loading the environment variables necessary to run the program.
RUN apt update
# TODO: for development purposes, remove
RUN apt install postgresql-client nmap net-tools -y

ENV PURITY_DB_HOST="localhost"
ENV PURITY_DB_PORT="5432"
ENV PURITY_DB_NAME="purity"
ENV PURITY_DB_USER="postgres"
ENV PURITY_DB_PASS="eGGQATPkx8JA66"
ENV PURITY_DB_SSL_MODE="disable"
ENV GOOGLE_APPLICATION_CREDENTIALS="/home/gb/Documents/purity-vision-5749de290743.json"
ENV PROJECT_ROOT="$(pwd)"
ENV PURITY_LOG_LEVEL="0"

WORKDIR /go/src/purity-vision
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["purity-vision-filter", "-port", "8080"]
