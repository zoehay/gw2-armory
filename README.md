# Guild Wars 2 Armoury 
This application allows GW2 players to view all of their character's inventories on one convenient, searchable page. 

## Technologies 
Go, PostgreSQL, GORM,
React, TypeScript, Vite

## Development
Run the backend with `go run cmd/main.go -v`. Run with set `APP_ENV=test` to enable service mocks.
Run the frontend with `npm run dev`.
Run backend tests with `go test ./...`

For local development with docker run `docker compose -f docker-compose-dev.yaml up --build`. Run with set `APP_ENV=test` to enable service mocks.
Run `npm run dev-build` from frontend to build with local docker dev baseURL.
