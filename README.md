# golden.rama.backendtemplate

## Start the API

API server [http://localhost:1323](http://localhost:1323)

```shell
export GOPRIVATE=github.com/Golden-Rama-Digital,github.com/harryosmar,github.com/tripdeals

go run main.go
```

## Start API using Docker

```shell
docker-compose down && \
docker volume rm paymentservice_mysql_data && \
docker-compose up -d
```

## Dependencies

- [github.com/Golden-Rama-Digital/library-core-go](https://github.com/Golden-Rama-Digital/library-core-go)
- [github.com/Golden-Rama-Digital/library-espay-go](https://github.com/Golden-Rama-Digital/library-espay-go)
