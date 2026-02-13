#!/bin/bash

docker network inspect "armory-network" > /dev/null
if [ $? -ne 0 ]; then
    docker network create -d bridge armory-network
else 
    echo "network already created"
fi

docker container create \
    --name armory-db \
    --network armory-network \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_DB=armory \
    -e POSTGRES_PASSWORD_FILE=/run/secrets/db_password.txt \
    -v armorydb:/var/lib/postgresql/data \
    -v $(pwd)/secrets/armory_db_password.txt:/run/secrets/db_password.txt:ro \
    postgres:16

docker container create \
    --name armory-backend \
    --restart unless-stopped \
    -e ARMORY_DB_PASSWORD_FILE=/run/secrets/db_password.txt \
    -e CORS_ALLOW_ORIGIN='https://localhost http://localhost' \
    --network armory-network \
    -v $(pwd)/secrets/armory_db_password.txt:/run/secrets/db_password.txt:ro \
    armory-backend
    
# nginx with mounted conf file for dev
docker container create \
    --name nginx \
    --restart unless-stopped \
    --network armory-network \
    -v $(pwd)/nginx/local-certs/certs:/etc/nginx/certs:ro \
    -v $(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
    -p 80:80 \
    -p 443:443 \
    armory-nginx