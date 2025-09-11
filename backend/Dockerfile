# syntax=docker/dockerfile:1

FROM golang:1.24 AS build

# Set destination for COPY
WORKDIR /src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/main.go

# Use light weight image for prod later 
FROM golang:1.24

WORKDIR /

COPY --from=build /backend /backend

EXPOSE 8000
ARG APP_ENV
ENV APP_ENV=$APP_ENV
ENV DB_DSN_FILE=/run/secrets/db_dsn

# Run
CMD ["/backend"]
