#!/usr/bin/env bash

db_name="purity-pg"

if [[ -z $(docker ps --filter=name=purity-pg -q) ]]; then
    echo "Starting database container"
    container=$(docker run --name ${db_name} \
	   -p "${PURITY_DB_PORT}":5432 \
	   -e POSTGRES_DB="${PURITY_DB_NAME}" \
	   -e POSTGRES_PASSWORD="${PURITY_DB_PASS}" \
	   -v /pg-docker-data:/var/lib/postgresql/data \
	   -v "${PROJECT_ROOT}/build":/docker-entrypoint-initdb.d \
	   --rm \
	   --detach \
	   postgres)
else
    echo "Database container already running"
    container=$(docker ps --filter=name=${db_name} -q)
fi

echo "Waiting for Postgres to be live"
while true
do
    echo "."
    if docker exec "${container}" pg_isready &> /dev/null; then
	break
    fi
done

