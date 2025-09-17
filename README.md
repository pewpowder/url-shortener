# url-shortener

## Preparing the app

1. You need to start DB locally and set up /internal/resources/config/app.development.yaml db section accordingly


And before starting run migration file from source directory

```
go run cmd/migrations/main.go

```

To start the application locally

```
go run cmd/url-shortener/main.go

```

with hot reload

```
air

```

