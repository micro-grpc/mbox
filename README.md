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

## Сборка Docker образа
 
Если вы используете приватный репозиторий то необходимо прописать доступ в файле 
который размещен в Вашей домашней директории ***.netrc*

```bash
machine github.com login [YOUR_GITHUB_USERNAME] password [YOUR_GITHUB_TOKEN]
machine gitlab.com login [YOUR_GITLAB_USERNAME] password [YOUR_GITLAB_TOKEN]

```

Токены вы можете сгенерировать на страницах [Personal GitHub Token](https://github.com/settings/tokens) [Personal GitLab Token](https://gitlab.com/profile/personal_access_tokens).
