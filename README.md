# Golang micro framework for gRPC microservice

[Validator](https://github.com/mwitkow/go-proto-validators)

[Recovery](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/recovery)

[Prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus)


# Для клиента

[Prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus)

[Client-Side Request Retry Interceptor](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/retry)


```bash
./ping client -vvvv --debug --data=pong --panic 1 --auth=basic --login=dev --passwd=disabled --repeat=2
./ping client -vvvv --debug --data=pong --panic 1 --auth=basic --login=dev --passwd=12345678 --repeat=2
./ping client -vvvv --debug --data=pong --panic 1 --auth=jwt --repeat=3
```