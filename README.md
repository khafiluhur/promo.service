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
docker volume rm promoservice_mysql_data && \
docker-compose up -d
```

## Dependencies

- [github.com/Golden-Rama-Digital/library-core-go](https://github.com/Golden-Rama-Digital/library-core-go)
- [github.com/tripdeals/cms.backend.tripdeals.id](https://github.com/tripdeals/cms.backend.tripdeals.id)
- [github.com/tripdeals/payment.service](https://github.com/tripdeals/payment.service)
- [github.com/tripdeals/library-service.go](https://github.com/tripdeals/library-service.go)
- [github.com/tripdeals/library-service.go](https://github.com/tripdeals/library-service.go)
- [github.com/harryosmar/cache-go](https://github.com/harryosmar/cache-go)
- [github.com/harryosmar/generic-gorm](https://github.com/harryosmar/generic-gorm)
