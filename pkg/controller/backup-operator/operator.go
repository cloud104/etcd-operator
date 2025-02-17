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

package controller

import (
	"context"
	"os"
	"sync"

	api "github.com/cloud104/etcd-operator/pkg/apis/etcd/v1beta2"
	"github.com/cloud104/etcd-operator/pkg/client"
	"github.com/cloud104/etcd-operator/pkg/generated/clientset/versioned"
	"github.com/cloud104/etcd-operator/pkg/util/constants"
	"github.com/cloud104/etcd-operator/pkg/util/k8sutil"

	"github.com/sirupsen/logrus"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Backup struct {
	logger *logrus.Entry

	namespace string
	// k8s workqueue pattern
	indexer  cache.Indexer
	informer cache.Controller
	queue    workqueue.RateLimitingInterface

	kubecli     kubernetes.Interface
	backupCRCli versioned.Interface
	kubeExtCli  apiextensionsclient.Interface

	backupRunnerStore sync.Map
}

type BackupRunner struct {
	spec       api.BackupSpec
	cancelFunc context.CancelFunc
}

// New creates a backup operator.
func New() *Backup {
	return &Backup{
		logger:      logrus.WithField("pkg", "controller"),
		namespace:   os.Getenv(constants.EnvOperatorPodNamespace),
		kubecli:     k8sutil.MustNewKubeClient(),
		backupCRCli: client.MustNewInCluster(),
		kubeExtCli:  k8sutil.MustNewKubeExtClient(),
	}
}

// Start starts the Backup operator.
func (b *Backup) Start(ctx context.Context) error {
	go b.run(ctx)
	<-ctx.Done()
	return ctx.Err()
}
