#!/bin/bash
 create() {
  kind create cluster --name etcd-operator --config=hack/developer/kind-cluster-config.yaml
  kind get kubeconfig --name etcd-operator > kubeconfig.yaml
  export KUBECONFIG=$pwd/develop/kubeconfig.yaml
  #kubectl taint nodes <node-name-medium> etcd-size=medium:NoSchedule
  #kubectl taint nodes <node-name-large> etcd-size=large:NoSchedule
}

 delete(){
  kind delete cluster --name etcd-operator
}