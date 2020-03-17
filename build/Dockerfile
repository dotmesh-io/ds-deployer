FROM golang:alpine AS build-env
WORKDIR /usr/local/go/src/github.com/dotmesh-io/ds-deployer
COPY . /usr/local/go/src/github.com/dotmesh-io/ds-deployer

ENV GO111MODULE=off

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make

RUN make install-release

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL name="Dotscience Model Deployment Service" \
      vendor="Dotscience." \
      version="v0.1.0" \
      release="1" \
      distribution-scope="public" \
      architecture="x86_64" \
      url="https://docs.dotscience.com" \
      io.k8s.description="Dotscience operators runs within a Openshift (or Kubernetes) cluster, and provides a means to deploy and manage ML models, configuring deployments, services and ingresses." \
      io.k8s.display-name="Dotscience Operator" \
      io.openshift.tags="dotscience,ml,operator,scm,git" \
      summary="Dotscience Model Deployment Service" \
      description="Dotscience Deployment operators runs within a Openshift (or Kubernetes) cluster, and provides a means to deploy and manage ML models, configuring deployments, services and ingresses."

ENV OPERATOR=/usr/local/bin/ds-deployer \
    USER_UID=1001 \
    USER_NAME=ds-deployer

# get licenses
COPY --from=build-env /usr/local/go/src/github.com/dotmesh-io/ds-deployer/licenses /licenses

COPY --from=build-env /usr/local/go/bin/ds-deployer ${OPERATOR}

# Copy endtrypoint and user_setup
COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
