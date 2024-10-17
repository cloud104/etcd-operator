#!/usr/bin/env bash

echo "generate certificates for etcd-cluster"
(bash ./gen-cert.sh)

echo "create secrets for etcd"
kubectl create -n default secret generic etcd-server-tls --from-file=server-ca.crt --from-file=server.crt --from-file=server.key
kubectl create -n default secret generic etcd-client-tls --from-file=etcd-client-ca.crt --from-file=etcd-client.crt --from-file=etcd-client.key
kubectl create -n default secret generic etcd-peer-tls --from-file=peer-ca.crt --from-file=peer.crt --from-file=peer.key