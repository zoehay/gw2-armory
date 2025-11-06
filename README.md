# Guild Wars 2 Armoury 
This application allows GW2 players to view all of their character's inventories on one convenient, searchable page. 

## Technologies 
Go, PostgreSQL, GORM,
React, TypeScript, Vite

## Development
Run the backend with `go run cmd/main.go -v`. 
Run the frontend with `npm run dev`.
Run backend tests with `go test ./...`

When running the backend, setting `APP_ENV=test` enables mocks in main.go and selects dsn for the testing database in router.go. 

For local development with docker run `docker compose -f docker-compose-dev.yaml up --build`. 

The nginx container will copy www into /var/www to serve. 

For local dev with nginx SSL configuration, from directory local-certs `mkdir certs` and run local-certs.sh for localhost cert using mkcert.

