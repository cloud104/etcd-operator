# Build step: docker build --tag gcr.io/coreos-k8s-scale-testing/logcollector -f hack/build/logcollector/Dockerfile .

FROM golang:1.23.3-alpine3.20

ADD ./ /go/src/github.com/cloud104/etcd-operator

WORKDIR /go/src/github.com/cloud104/etcd-operator

RUN rm -rf _output _test .git .gitignore

RUN go build -o /usr/local/bin/logcollector test/logcollector/main.go
