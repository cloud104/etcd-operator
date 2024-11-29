#!/usr/bin/env bash

kind=$(which tcloud-cli)
kubectl=$(which kubectl)
skaffold=$(which skaffold)

# Verifica se o kind está instalado
if [ -z "$kind" ]; then
    echo "O tcloud-cli não está instalado. Instale-o primeiro."
    [ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.24.0/kind-$(uname)-amd64
    sudo install $kind /usr/local/bin/kind

fi
# Verifica se o tcloud-cli está instalado
if [ -z "$kubectl" ]; then
   echo "O kubectl não está instalado. Instale-o primeiro."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    sudo install kubectl /usr/local/bin/
fi
# Verifica se o skaffold está instalado
if [ -z "$skaffold" ]; then
   echo "O skaffold não está instalado. Instale-o primeiro."
   # For Linux x86_64 (amd64)
   curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
   sudo install skaffold /usr/local/bin/
fi

 create() {
  kind create cluster --name etcd-operator --config=hack/developer/kind-cluster-config.yaml
  kind get kubeconfig --name etcd-operator > $PWD/hack/developer/kubeconfig.yaml
  export KUBECONFIG=$PWD/hack/developer/kubeconfig.yaml
  kubectl apply -f ${PWD}/example/crd/
  (cd ${PWD}/example/rbac && bash create_role.sh --namespace=default)
#  kubectl apply -f ${PWD}/example/crd/
#  kubectl apply -f ${PWD}/example/deployment.yaml
}

 delete(){
  kind delete cluster --name etcd-operator
  rm -rf $PWD/hack/developer/kubeconfig.yaml
}

$1