JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -s -w
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.BuildDate=$(JOBDATE)

install-release:
	@echo "++ Installing Dotscience Deployer"	
	go install -ldflags="$(LDFLAGS)" github.com/dotmesh-io/ds-deployer/cmd/ds-deployer

install:
	cd cmd/ds-deployer && go install

proto:	
	cd apis/deployer/v1 && protoc --gofast_out=plugins=grpc:. deployer.proto

image:
	docker build -t quay.io/dotmesh/dotscience-deployer:alpha -f Dockerfile .
	docker push quay.io/dotmesh/dotscience-deployer:alpha

run: install
	ds-deployer run --no-incluster --no-require-tls

test:
	go get github.com/mfridman/tparse
	go test -json -v `go list ./... | egrep -v /tests` -cover | tparse -all -smallscreen