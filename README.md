# mBox golang micro framework for gRPC microservice

## Installation

Installing **mBox**, you can install the cmd line app to generate new micro services and the required libraries. First you'll need Google's Protocol Buffers installed.

```bash
brew install protobuf
go get -u github.com/micro-grpc/mbox/...
sudo mbox --bash-completion
```

## Getting Started

To generate a new service, run mBox in new trminal with a short folder path.

```bash
mbox init github.com/micro-grpc/example-service --name=example
```


[Validator](https://github.com/mwitkow/go-proto-validators)

[Recovery](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/recovery)

[Prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus)


## Для клиента

[Prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus)

[Client-Side Request Retry Interceptor](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/retry)


```bash
./ping client -vvvv --debug --data=pong --panic 1 --auth=basic --login=dev --passwd=disabled --repeat=2
./ping client -vvvv --debug --data=pong --panic 1 --auth=basic --login=dev --passwd=12345678 --repeat=2
./ping client -vvvv --debug --data=pong --panic 1 --auth=jwt --repeat=3
```