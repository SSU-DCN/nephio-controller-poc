FROM ubuntu:18.04
MAINTAINER "ntnguyencse@gmail.com"
ARG GO_VERSION
ENV GO_VERSION=1.18.3

RUN apt-get update && apt-get install -y wget git gcc curl file && wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz" && rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR /workspace
RUN wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 && chmod a+x /usr/local/bin/yq
RUN wget  https://dl.k8s.io/release/v1.25.4/bin/linux/amd64/kubectl && chmod +x kubectl && mkdir -p ~/.local/bin && cp ./kubectl ~/.local/bin/kubectl
RUN wget  https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.3.0/clusterctl-linux-amd64 && mv clusterctl-linux-amd64 clusterctl && chmod +x clusterctl && cp ./clusterctl ~/.local/bin/clusterctl
ENV PATH ~/.local/bin:$PATH

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile
COPY run.sh run.sh
# COPY env.rc env.rc
# COPY clouds.yaml clouds.yaml
COPY config-env config-env
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
ENV KUBECONFIG /kubeconfig/config
RUN chmod +x run.sh
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o provider main.go

ENTRYPOINT  ["./run.sh"]
#WORKDIR $G
