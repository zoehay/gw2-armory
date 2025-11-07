

# syntax=docker/dockerfile:1

FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release

WORKDIR /

COPY --from=build /backend /backend

EXPOSE 8000

ARG APP_ENV
ENV APP_ENV=$APP_ENV
ARG GIN_MODE
ENV GIN_MODE=$GIN_MODE

USER nonroot:nonroot

CMD ["/backend"]

