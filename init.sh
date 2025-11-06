#!/bin/bash

docker network inspect "armory-network" > /dev/null
if [ $? -ne 0 ]; then
    docker network create -d bridge armory-network
else 
    echo "network already created"
fi

sh db/init-db.sh && \

sh backend/init-backend.sh