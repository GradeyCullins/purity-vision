#!/usr/bin/env bash

if [[ -z $(docker ps --filter=name=purity-pg -q) ]]; then
    echo "Starting database container"
    docker run --name purity-pg \
	   -p ${PURITY_DB_PORT}:5432 \
	   -e POSTGRES_DB="${PURITY_DB_NAME}" \
	   -e POSTGRES_PASSWORD="${PURITY_DB_PASS}" \
	   -v /pg-docker-data:/var/lib/postgresql/data \
	   -v ${PROJECT_ROOT}/build:/docker-entrypoint-initdb.d \
	   --rm \
	   --detach \
	   postgres
else
    echo "Database container already running"
fi

