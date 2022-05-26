/*
Copyright 2022 The etcd-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta2

import (
	v1beta2 "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EtcdRestoreLister helps list EtcdRestores.
// All objects returned here must be treated as read-only.
type EtcdRestoreLister interface {
	// List lists all EtcdRestores in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.EtcdRestore, err error)
	// EtcdRestores returns an object that can list and get EtcdRestores.
	EtcdRestores(namespace string) EtcdRestoreNamespaceLister
	EtcdRestoreListerExpansion
}

// etcdRestoreLister implements the EtcdRestoreLister interface.
type etcdRestoreLister struct {
	indexer cache.Indexer
}

// NewEtcdRestoreLister returns a new EtcdRestoreLister.
func NewEtcdRestoreLister(indexer cache.Indexer) EtcdRestoreLister {
	return &etcdRestoreLister{indexer: indexer}
}

// List lists all EtcdRestores in the indexer.
func (s *etcdRestoreLister) List(selector labels.Selector) (ret []*v1beta2.EtcdRestore, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.EtcdRestore))
	})
	return ret, err
}

// EtcdRestores returns an object that can list and get EtcdRestores.
func (s *etcdRestoreLister) EtcdRestores(namespace string) EtcdRestoreNamespaceLister {
	return etcdRestoreNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EtcdRestoreNamespaceLister helps list and get EtcdRestores.
// All objects returned here must be treated as read-only.
type EtcdRestoreNamespaceLister interface {
	// List lists all EtcdRestores in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.EtcdRestore, err error)
	// Get retrieves the EtcdRestore from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.EtcdRestore, error)
	EtcdRestoreNamespaceListerExpansion
}

// etcdRestoreNamespaceLister implements the EtcdRestoreNamespaceLister
// interface.
type etcdRestoreNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EtcdRestores in the indexer for a given namespace.
func (s etcdRestoreNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.EtcdRestore, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.EtcdRestore))
	})
	return ret, err
}

// Get retrieves the EtcdRestore from the indexer for a given namespace and name.
func (s etcdRestoreNamespaceLister) Get(name string) (*v1beta2.EtcdRestore, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("etcdrestore"), name)
	}
	return obj.(*v1beta2.EtcdRestore), nil
}
