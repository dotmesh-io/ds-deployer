JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -s -w
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/dotmesh-io/ds-deployer/pkg/version.BuildDate=$(JOBDATE)

install-release:
	@echo "++ Installing Dotscience Deployer"	
	CGO_ENABLED=0 go install -ldflags="$(LDFLAGS)" github.com/dotmesh-io/ds-deployer/cmd/ds-deployer

install:
	cd cmd/ds-deployer && go install

gen:
	cd apis/deployer/v1 && protoc --gofast_out=plugins=grpc:. deployer.proto
	easyjson pkg/health/health_server.go

image:
	docker build -t quay.io/dotmesh/dotscience-deployer:alpha -f Dockerfile .
	
image-push: 
	docker push quay.io/dotmesh/dotscience-deployer:alpha

# this scan.connect.redhat.com registry is where we need to push images for scanning. Once it's scanned, you can "publish"
# image here https://connect.redhat.com/project/2420521/view. To get the push
image-ubi:
	docker --config ~/.rh/ build -t scan.connect.redhat.com/ospid-54a873aa-9b87-4659-9f5e-f4632074a5a7/dotscience-deployer-ubi7:beta-rc2 -f Dockerfile.ubi8 .

# Login here using their generated token. To get the push config, go here
# https://connect.redhat.com/project/2420521/view and click "Upload Your Image", they will
# give you a token with incorrect login instructions :D. Just do:
# docker login -u unused -p TOKEN_HERE scan.connect.redhat.com
image-push-ubi: image-ubi
	docker push scan.connect.redhat.com/ospid-54a873aa-9b87-4659-9f5e-f4632074a5a7/dotscience-deployer-ubi7:beta-rc2

run: install
	ds-deployer run --no-incluster --no-require-tls

test:
	go get github.com/mfridman/tparse
	go test -json -v `go list ./... | egrep -v /tests` -cover | tparse -all -smallscreen