// Copyright 2017 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2eutil

import (
	"context"
	api "github.com/cloud104/etcd-operator/pkg/apis/etcd/v1beta2"
	"github.com/cloud104/etcd-operator/pkg/generated/clientset/versioned"
	"github.com/cloud104/etcd-operator/pkg/util/k8sutil"
	"github.com/cloud104/etcd-operator/pkg/util/retryutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StorageCheckerOptions struct {
	S3Cli          *s3.S3
	S3Bucket       string
	DeletedFromAPI bool
}

func CreateCluster(t *testing.T, crClient versioned.Interface, namespace string, cl *api.EtcdCluster) (*api.EtcdCluster, error) {
	ctx := context.Background()
	cl.Namespace = namespace
	res, err := crClient.EtcdV1beta2().EtcdClusters(namespace).Create(ctx, cl, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	t.Logf("creating etcd cluster: %s", res.Name)

	return res, nil
}

func UpdateCluster(crClient versioned.Interface, cl *api.EtcdCluster, maxRetries int, updateFunc k8sutil.EtcdClusterCRUpdateFunc) (*api.EtcdCluster, error) {
	return AtomicUpdateClusterCR(crClient, cl.Name, cl.Namespace, maxRetries, updateFunc)
}

func AtomicUpdateClusterCR(crClient versioned.Interface, name, namespace string, maxRetries int, updateFunc k8sutil.EtcdClusterCRUpdateFunc) (*api.EtcdCluster, error) {
	ctx := context.Background()
	result := &api.EtcdCluster{}
	err := retryutil.Retry(1*time.Second, maxRetries, func() (done bool, err error) {
		etcdCluster, err := crClient.EtcdV1beta2().EtcdClusters(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		updateFunc(etcdCluster)

		result, err = crClient.EtcdV1beta2().EtcdClusters(namespace).Update(ctx, etcdCluster, metav1.UpdateOptions{})
		if err != nil {
			if apierrors.IsConflict(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return result, err
}

func DeleteCluster(t *testing.T, crClient versioned.Interface, kubeClient kubernetes.Interface, cl *api.EtcdCluster) error {
	ctx := context.Background()
	t.Logf("deleting etcd cluster: %v", cl.Name)
	err := crClient.EtcdV1beta2().EtcdClusters(cl.Namespace).Delete(ctx, cl.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return waitResourcesDeleted(t, kubeClient, cl)
}
