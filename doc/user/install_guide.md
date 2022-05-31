# Installation guide

## Set up RBAC

Set up basic [RBAC rules][rbac-rules] for etcd operator:

```bash
$ example/rbac/create_role.sh
```
## Install CRD
```bash
$ kubectl create -f example/crd
```

## Install etcd operator

Create a deployment for etcd operator:

```bash
$ kubectl create -f example/deployment.yaml
```

## Uninstall etcd operator

Note that the etcd clusters managed by etcd operator will **NOT** be deleted even if the operator is uninstalled.

This is an intentional design to prevent accidental operator failure from killing all the etcd clusters.

To delete all clusters, delete all cluster CR objects before uninstalling the operator.

Clean up etcd operator:

```bash
kubectl delete -f example/deployment.yaml
kubectl delete endpoints etcd-operator
kubectl delete -f example/crd
kubectl delete clusterrole etcd-operator
kubectl delete clusterrolebinding etcd-operator
```

[rbac-rules]: rbac.md