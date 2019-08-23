
install:
	cd cmd/ds-deployer && go install

proto:
	cd apis/deployer/v1 && protoc --gofast_out=. deployer.proto

image:
	# dotmesh/dotscience-deployer:latest
	docker build -t quay.io/dotmesh/dotscience-deployer:alpha -f Dockerfile .
	docker push quay.io/dotmesh/dotscience-deployer:alpha