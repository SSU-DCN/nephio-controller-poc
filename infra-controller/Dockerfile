FROM ubuntu:18.04
MAINTAINER "ntnguyencse@gmail.com"
ARG GO_VERSION
ENV GO_VERSION=1.18.3

RUN apt-get update && apt-get install -y wget git gcc && wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz" && rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
RUN wget  https://dl.k8s.io/release/v1.25.4/bin/linux/amd64/kubectl && chmod +x kubectl && mkdir -p ~/.local/bin && mv ./kubectl ~/.local/bin/kubectl
ENV PATH ~/.local/bin:$PATH

ENTRYPOINT  ["/bin/bash"]
#WORKDIR $G
