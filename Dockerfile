FROM golang:alpine AS build-env
WORKDIR /usr/local/go/src/github.com/dotmesh-io/ds-deployer
COPY . /usr/local/go/src/github.com/dotmesh-io/ds-deployer

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make

ENV GO111MODULE=off

RUN make install-release

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build-env /usr/local/go/bin/ds-deployer /bin/ds-deployer
CMD ["ds-deployer"]
