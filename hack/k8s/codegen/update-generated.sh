#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

bash -ax ~/go/pkg/mod/k8s.io/code-generator@v0.24.1/generate-groups.sh \
  "all" \
  "github.com/cloud104/etcd-operator/pkg/generated" \
  "github.com/cloud104/etcd-operator/pkg/apis" \
  "etcd:v1beta2" \
  --go-header-file "./hack/k8s/codegen/boilerplate.go.txt" \
  $@
