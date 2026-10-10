# Guild Wars 2 Armory

![Frontend CI](https://github.com/zoehay/gw2-armory/actions/workflows/frontend.yml/badge.svg)

A web application that lets Guild Wars 2 players view the inventories of all their characters, their bank, and their material storage on a single searchable page.

Players sign in with a [Guild Wars 2 API key](https://account.arena.net/applications). The backend fetches account and character data from the official GW2 API, stores it in PostgreSQL, and serves it to a React frontend for browsing and filtering.

## Features

- **Unified inventory view.** See every character's bags, the account bank, and material storage together.
- **Search and filtering.** Find an item without logging in to each character.
- **API key management.** Add and remove GW2 API keys from the app.
- **Guest sessions.** Try the app with an API key without creating an account.

## Tech Stack

| Layer    | Technologies                                                      |
| -------- | ----------------------------------------------------------------- |
| Backend  | Go, [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/)       |
| Database | PostgreSQL                                                        |
| Frontend | React, TypeScript, Vite, React Router                             |
| Testing  | Go `testing`, Vitest, React Testing Library                       |
| Infra    | Docker Compose, nginx (TLS termination and reverse proxy), GitHub Actions |

## Architecture

```
            ┌──────────────────────────────┐
 Browser ──►│ nginx (:80 / :443)           │
            │  /       → static frontend   │
            │  /api/*  → backend :8000     │
            └──────────────┬───────────────┘
                           │
                 ┌─────────▼─────────┐        ┌──────────────────┐
                 │ Go / Gin backend  │───────►│ GW2 public API   │
                 └─────────┬─────────┘        └──────────────────┘
                           │
                 ┌─────────▼─────────┐
                 │ PostgreSQL        │
                 └───────────────────┘
```

## Project Structure

```
.
├── backend/
│   ├── cmd/main.go            # Entry point
│   ├── internal/
│   │   ├── api/               # Handlers, middleware (CORS, sessions), routes
│   │   ├── db/                # Database setup, models, repositories
│   │   ├── gw2_client/        # GW2 API client, providers, and mocks
│   │   └── services/          # Business logic
│   ├── test_data/             # Fixture responses from the GW2 API
│   └── tests/                 # Integration and handler tests
├── frontend/                  # React + TypeScript single-page app
├── nginx/                     # nginx images, configs, and local cert script
├── db/                        # Database container scripts
├── docker-compose-dev.yaml
└── docker-compose-prod.yaml
```

## Getting Started

### Prerequisites

- Go 1.27+
- Node.js 24+ and npm
- PostgreSQL (or Docker)
- [mkcert](https://github.com/FiloSottile/mkcert), for local HTTPS through nginx
- A Guild Wars 2 API key with the **account**, **inventories**, and **characters** permissions

### Configuration

Environment files are gitignored. Copy the examples and fill in your own values:

```sh
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

**`backend/.env`**

| Variable             | Description                                                                       |
| -------------------- | --------------------------------------------------------------------------------- |
| `DEV_NO_MOCK_DB_DSN` | PostgreSQL DSN for local development                                              |
| `TEST_DB_DSN`        | PostgreSQL DSN for the test database                                              |
| `CORS_ALLOW_ORIGIN`  | Space-separated allowed origins. Defaults to `https://localhost:5173 http://localhost:5173` |
| `DOMAIN`             | Domain for session cookies. Defaults to `localhost`                               |
| `APP_ENV`            | Set to `test` to use mocked GW2 API providers and the test database               |

When `ARMORY_DB_PASSWORD_FILE` is set, as it is in Docker, the backend reads the database password from that file and ignores `.env`.

**`frontend/.env`**

| Variable           | Description              |
| ------------------ | ------------------------ |
| `VITE_APP_API_URL` | Base URL of the backend API |

### Running Locally

Start the backend (listens on port `8000`):

```sh
cd backend
go run cmd/main.go
```

Start the frontend dev server (port `5173`):

```sh
cd frontend
npm install
npm run dev
```

Set `APP_ENV=test` when starting the backend to replace calls to the GW2 API with mocks. This is useful for working offline or without an API key.

### Running with Docker

The development stack runs PostgreSQL, the backend, and nginx with local TLS.

1. Create the database password secret at `secrets/armory_db_password.txt`.
2. Generate a localhost certificate with mkcert:

   ```sh
   cd nginx/local-certs
   mkdir -p certs
   sh local-certs.sh
   ```
3. Build the frontend and place the output in `www/armory/`. The nginx image copies `www/` into `/var/www` at build time.
4. Start the stack:

   ```sh
   docker compose -f docker-compose-dev.yaml up --build
   ```

The app is then served at `https://localhost`, and requests to `/api/*` are proxied to the backend.

## Testing

### Backend

```sh
cd backend
go test -p 1 -v ./tests/...
```

Test packages share one database, so they must run serially (`-p 1`). Add the following filter for condensed output:

```sh
go test -p 1 -v ./tests/... 2>&1 | grep -E "^(--- FAIL|--- PASS|FAIL|ok)"
```

### Frontend

```sh
cd frontend
npm test          # run once
npm run test:watch
npm run lint
npm run format:check
```

GitHub Actions runs linting, formatting checks, tests, and a production build on every push and pull request that touches `frontend/`.

## API Overview

| Method | Endpoint                                         | Description                         |
| ------ | ------------------------------------------------ | ----------------------------------- |
| GET    | `/items`, `/items/:id`                           | Item metadata                       |
| POST   | `/apikeys`                                       | Register an API key and start a session |
| POST   | `/login`                                         | Log in                              |
| GET    | `/account/info`                                  | Current account details             |
| GET    | `/account/inventory`                             | All items across the account        |
| GET    | `/account/accountinventory`                      | Bank and material storage           |
| GET    | `/account/characters/:charactername/inventory`   | One character's inventory           |
| POST   | `/account/searchinventory`                       | Filtered inventory search           |
| DELETE | `/account/delete`                                | Delete the account and its data     |

All `/account` routes require a valid session cookie.

## Acknowledgements

Guild Wars 2 is a trademark of ArenaNet, LLC. This project is not affiliated with or endorsed by ArenaNet. Game data comes from the [official Guild Wars 2 API](https://wiki.guildwars2.com/wiki/API:Main).
