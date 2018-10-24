#
# Makefile for this application
#
-include variable.mak

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


REPO_URL ?=
IMAGE_NAME ?=
USER_NAME ?= grengojbo
ADMIN_USER ?= grengojbo
TAG_VERSION=$(shell cat RELEASE)

OSNAME=$(shell uname)
GO=$(shell which go)

GO_MODULE=on

CUR_TIME=$(shell date '+%Y-%m-%d_%H:%M:%S')
# Program version
VERSION=$(shell cat RELEASE)

# Binary name for bintray
BIN_NAME=$(shell basename $(abspath ./))

# Project name for bintray
PROJECT_NAME=$(shell basename $(abspath ./))
PROJECT_DIR=$(shell pwd)

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com.org

# Grab the current commit
GIT_COMMIT="$(shell git rev-parse HEAD)"

# Check if there are uncommited changes
GIT_DIRTY="$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

DIST_DIR="${PROJECT_DIR}/dist"
DIST_BIN="${DIST_DIR}/bin"

NAME ?= ${PROJECT_NAME}
TAG=${USER_NAME}/$(PROJECT_NAME)$(IMAGE_NAME):$(TAG_VERSION)

BUILD_TAGS ?= "netgo"
BUILD_TAGS_BINDATA ?= "netgo bindatafs"
BUILD_ENV = GOOS=linux GOARCH=amd64
ENVFLAGS = CGO_ENABLED=1 $(BUILD_ENV)
ifneq ($(GOOS), darwin)
  EXTLDFLAGS = -extldflags "-lm -lstdc++ -static"
else
  EXTLDFLAGS =
endif

GO_LINKER_FLAGS ?= -ldflags '$(EXTLDFLAGS) -s -w \
  -X "main.BuildTime=${CUR_TIME}" \
  -X "main.Version=${VERSION}" \
  -X "main.GitHash=${GIT_COMMIT}" \
  -X "config.Version=${VERSION}"

# Add the godep path to the GOPATH
#GOPATH=$(shell godep path):$(shell echo $$GOPATH)

#ifeq ($(OS),Darwin)
#  URL=$(shell dinghy ip)
#else
#  URL="127.0.0.1"
#endif
URL=$(shell dinghy ip)

default: help

help:
	@echo "..............................................................."
	@echo "Project: $(PROJECT_NAME) | current dir: $(PROJECT_DIR)"
	@echo "version: $(VERSION) GIT_DIRTY: $(GIT_DIRTY)\n"
	@echo "make init        - Load project"
	@echo "make protoc      - Generate gRPC"
	@echo "make build       - Build for current OS project"
	@echo "make release     - Build release project"
	@echo "make docs        - Project documentation"
	@echo "make deploy      - Deploy bin files to server"
	@echo "make update      - Update vendor files"
	@echo "make test        - Run all test"
	@echo "make version     - Current project version"
	@echo "make serve       - Run local Docker image"
	@echo "make build-docker - Build Docker image"
	@echo "make push         - Push Docker image"
	@echo "make publish      - Publication new Release"
	@echo "...............................................................\n"

init:
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/vgo
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/proto
	@go get -u github.com/golang/protobuf/ptypes
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

update:
	@#dep ensure -update
	@GO111MODULE=${GO_MODULE} go get -u ./...

protoc:
	@echo "Generate gRPC"
	@protoc --proto_path=${GOPATH}/src \
	 --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	 -I types/ \
	 --go_out=plugins=grpc:types \
	 --govalidators_out=types \
	 types/ping/*.proto
	@protoc --proto_path=${GOPATH}/src \
	 --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	 -I . \
	 --go_out=. \
	 ./grpc/pb/authorize/*.proto


deploy:
	@echo "TODO"

publish:
	@#./bumper.sh
	@git add -A
	@git commit -am "Bump version to v$(shell cat RELEASE)"
	@git tag v$(shell cat RELEASE)
	@git push
	@git push --tags

release: clean
	@mkdir -p $(DIST_BIN)
	@echo "building release ${BIN_NAME} ${VERSION}"
	@GOOS=linux GOARCH=amd64 go build -a -tags "netgo bindatafs" -ldflags '-w -X main.BuildTime=${CUR_TIME} -X main.Version=${VERSION} -X main.GitHash=${GIT_COMMIT} -X config.Version=${VERSION}' -o $(DIST_HOME)/$(BIN_NAME) main.go
	@chmod 0755 $(DIST_BIN)/$(BIN_NAME)

clean:
	@test ! -e ./${BIN_NAME} || rm ./${BIN_NAME}
	@test ! -e ${DIST_HOME}/${BIN_NAME} || rm ${DIST_HOME}/${BIN_NAME}
	@#git gc --prune=0 --aggressive
	@find . -name "*.orig" -type f -delete
	@find . -name "*.log" -type f -delete
	@test ! -e ./dist || rm -R ./dist

test:
	@echo "Start test..."

clean-docker:
	docker rmi -f $(REPO_URL)$(TAG)
	docker system prune -f

push:
	docker push $(REPO_URL)$(TAG)

serve:
	@echo "$(OK_COLOR)RUN command line: open http://$(URL):$(PUB_PORT)/$(NO_COLOR)\n\n"
	@docker run --rm \
		--name=$(NAME) \
     	-p=$(PUB_PORT):$(PORT) \
     	-e PORT=$(PORT) \
     	-it $(REPO_URL)$(TAG)

stop-docker:
	docker stop $(NAME)

build-docker:
	docker build --tag=$(REPO_URL)$(TAG) .

# Attach a root terminal to an already running dev shell
shell:
	docker run -it --rm $(REPO_URL)$(TAG) bash

build-linux: clean
	@mkdir -p $(DIST_BIN)
	@echo "building version: ${VERSION} to  ${DIST_HOME}/${BIN_NAME}"
	@GO111MODULE=${GO_MODULE} GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -tags netgo -ldflags "-s -w -X main.release=${VERSION} -X main.Commit=${GIT_COMMIT} -X main.BuildTime=${CUR_TIME}" -o $(DIST_BIN)/$(BIN_NAME) main.go
	@echo " "

build: clean
	@mkdir -p $(DIST_BIN)
	@echo "building version: ${VERSION} to  ${DIST_HOME}/${BIN_NAME}"
	@GO111MODULE=${GO_MODULE} CGO_ENABLED=0 go build -a -installsuffix cgo -tags netgo -ldflags "-s -w -X main.release=${VERSION} -X main.Commit=${GIT_COMMIT} -X main.BuildTime=${CUR_TIME}" -o ./$(BIN_NAME) main.go
	@echo " "

version:
	@echo ${VERSION}

docs:
	godoc -http=:6060 -index

