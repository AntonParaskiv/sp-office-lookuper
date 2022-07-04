# sp-office-lookuper

### Run service

```shell
$ go run ./cmd/api
```

### Environments
```shell
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
HTTP_MAX_CONNECTIONS=1000
HTTP_PPROF_PORT=8180

GRPC_HOST=0.0.0.0
GRPC_PORT=8090

LOGGER_LEVEL=4 # info
LOGGER_DESTINATION=stdout|gelf
LOGGER_HOST=0.0.0.0 # for gelf only
LOGGER_PORT=12201 # for gelf only

JAEGER_HOST=localhost
JAEGER_PORT=6831
```

### Ports
* HTTP 8080
* HTTP PPROF 8180
* GRPC 8090