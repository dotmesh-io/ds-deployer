
install:
	cd cmd/ds-deployer && go install

proto:
	cd apis/deployer/v1 && protoc --gofast_out=. deployer.proto
