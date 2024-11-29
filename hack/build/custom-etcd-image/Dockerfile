FROM alpine:3.20.3 AS builder

ARG VERSION=$VERSION
WORKDIR /etcd
RUN apk update
RUN apk add curl
RUN curl -sL https://github.com/etcd-io/etcd/releases/download/v${VERSION}/etcd-v${VERSION}-linux-amd64.tar.gz -o etcd-v${VERSION}-linux-amd64.tar.gz && \
    tar -xzf etcd-v${VERSION}-linux-amd64.tar.gz -C /etcd/


# TODO: move to k8s.gcr.io/build-image/debian-base:bullseye-v1.y.z when patched
FROM debian:bullseye-20241016
#FROM debian:bullseye-20220328
ARG VERSION=$VERSION
#COPY --from=builder --chown=65532:65532 /workspace/etcd-operator /usr/local/bin/etcd-operator
COPY --from=builder /etcd/etcd-v${VERSION}-linux-amd64/etcd /usr/local/bin/
COPY --from=builder /etcd/etcd-v${VERSION}-linux-amd64/etcdctl /usr/local/bin/
COPY --from=builder /etcd/etcd-v${VERSION}-linux-amd64/etcdutl /usr/local/bin/
RUN mkdir -p /var/etcd/
RUN mkdir -p /var/lib/etcd/

# Alpine Linux doesn't use pam, which means that there is no /etc/nsswitch.conf,
# but Golang relies on /etc/nsswitch.conf to check the order of DNS resolving
# (see https://github.com/golang/go/commit/9dee7771f561cf6aee081c0af6658cc81fac3918)
# To fix this we just create /etc/nsswitch.conf and add the following line:
RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

EXPOSE 2379 2380 2381

# Define default command.
CMD ["/usr/local/bin/etcd"]