#/usr/bin/env bash

echo "Creating database container"
docker run --name purity-pg \
    -p ${PURITY_DB_PORT}:5432 \
    -e POSTGRES_DB="${PURITY_DB_NAME}" \
    -e POSTGRES_PASSWORD="${PURITY_DB_PASS}" \
    -v /pg-docker-data:/var/lib/postgresql/data \
    -v ${PROJECT_ROOT}/build:/docker-entrypoint-initdb.d \
    --rm \
    --detach \
    postgres