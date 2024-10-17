// Copyright 2016 The etcd-operator Authors
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

package cluster

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cloud104/etcd-operator/pkg/util/etcdutil"
	"github.com/cloud104/etcd-operator/pkg/util/k8sutil"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

func (c *Cluster) upgradeOneMember(memberName string) error {
	c.status.SetUpgradingCondition(c.cluster.Spec.Version)

	ns := c.cluster.Namespace

	pod, err := c.config.KubeCli.CoreV1().Pods(ns).Get(context.Background(), memberName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("fail to get pod (%s): %v", memberName, err)
	}
	oldpod := pod.DeepCopy()

	c.logger.Infof("upgrading the etcd member %v from %s to %s", memberName, k8sutil.GetEtcdVersion(pod), c.cluster.Spec.Version)
	pod.Spec.Containers[0].Image = k8sutil.ImageName(c.cluster.Spec.Repository, c.cluster.Spec.Version)
	clusterSpecVersion, err := semver.NewVersion(c.cluster.Spec.Version)
	if err != nil {
		return fmt.Errorf("erros parsing etcd version %s", c.cluster.Spec.Version)
	}
	switch {
	case clusterSpecVersion.LessThanEqual(semver.MustParse("3.5.6")):
		c.logger.Infof("pod: %v ", memberName)
		stringPodEtcdVersion := strings.Split(oldpod.Spec.Containers[0].Image, ":v")[1] //"quay.io/coreos/etcd" "3.5.15"
		podVersion, err := semver.NewVersion(stringPodEtcdVersion)
		if err != nil {
			return fmt.Errorf("parsing pod version: %v", err)
		}
		c.logger.Infof("podImage: %v ", podVersion)
		c.logger.Infof("pod running etcd version: %v", podVersion)
		if podVersion.Equal(semver.MustParse("3.5.7")) {
			err = c.rolloutOneMember(memberName, pod.Namespace)
			if err != nil {
				return fmt.Errorf("failed to rolloutOneMember: %v", err)
			}
		}
		c.logger.Infof("STATUS is: %v", c.cluster.Status.Members.Ready)
		k8sutil.SetEtcdVersion(pod, c.cluster.Spec.Version)
		patchdata, err := k8sutil.CreatePatch(oldpod, pod, v1.Pod{})
		if err != nil {
			return fmt.Errorf("error creating patch: %v", err)
		}
		_, err = c.config.KubeCli.CoreV1().Pods(ns).Patch(context.Background(), pod.GetName(), types.StrategicMergePatchType, patchdata, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("fail to update the etcd member (%s): %v", memberName, err)
		}
		c.logger.Infof("finished upgrading the etcd member %v", memberName)
		_, err = c.eventsCli.Create(context.Background(), k8sutil.MemberUpgradedEvent(memberName, k8sutil.GetEtcdVersion(oldpod), c.cluster.Spec.Version, c.cluster), metav1.CreateOptions{})
		if err != nil {
			c.logger.Errorf("failed to create member upgraded event: %v", err)
		}
		return nil
	case clusterSpecVersion.GreaterThanEqual(semver.MustParse("3.5.7")) && clusterSpecVersion.LessThanEqual(semver.MustParse("3.5.11")):
		c.logger.Infof("pod: %v ", memberName)
		stringPodEtcdVersion := strings.Split(oldpod.Spec.Containers[0].Image, ":v")[1] //"quay.io/coreos/etcd" "3.5.15
		podVersion, err := semver.NewVersion(stringPodEtcdVersion)
		if err != nil {
			return fmt.Errorf("parsing pod version: %v", err)
		}
		c.logger.Infof("podImage: %v ", podVersion)
		c.logger.Infof("pod running etcd version: %v", podVersion)
		if podVersion.Equal(semver.MustParse("3.5.6")) {
			err = c.rolloutOneMember(memberName, pod.Namespace)
			if err != nil {
				return fmt.Errorf("failed to rolloutOneMember: %v", err)
			}
		}
		if podVersion.Equal(semver.MustParse("3.5.12")) {
			err = c.rolloutOneMember(memberName, pod.Namespace)
			if err != nil {
				return fmt.Errorf("failed to rolloutOneMember: %v", err)
			}
		}
		c.logger.Infof("STATUS is: %v", c.cluster.Status.Members.Ready)
		k8sutil.SetEtcdVersion(pod, c.cluster.Spec.Version)
		patchdata, err := k8sutil.CreatePatch(oldpod, pod, v1.Pod{})
		if err != nil {
			return fmt.Errorf("error creating patch: %v", err)
		}
		_, err = c.config.KubeCli.CoreV1().Pods(ns).Patch(context.Background(), pod.GetName(), types.StrategicMergePatchType, patchdata, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("fail to update the etcd member (%s): %v", memberName, err)
		}
		c.logger.Infof("finished upgrading the etcd member %v", memberName)
		_, err = c.eventsCli.Create(context.Background(), k8sutil.MemberUpgradedEvent(memberName, k8sutil.GetEtcdVersion(oldpod), c.cluster.Spec.Version, c.cluster), metav1.CreateOptions{})
		if err != nil {
			c.logger.Errorf("failed to create member upgraded event: %v", err)
		}
		return nil
	case clusterSpecVersion.GreaterThanEqual(semver.MustParse("3.5.12")):
		c.logger.Infof("pod: %v ", memberName)
		stringPodEtcdVersion := strings.Split(oldpod.Spec.Containers[0].Image, ":v")[1] //"quay.io/coreos/etcd" "3.5.15
		podVersion, err := semver.NewVersion(stringPodEtcdVersion)
		if err != nil {
			return fmt.Errorf("parsing pod version: %v", err)
		}
		c.logger.Infof("podImage: %v ", podVersion)
		c.logger.Infof("pod running etcd version: %v", podVersion)
		if podVersion.Equal(semver.MustParse("3.5.11")) {
			err = c.rolloutOneMember(memberName, pod.Namespace)
			if err != nil {
				return fmt.Errorf("failed to rolloutOneMember: %v", err)
			}
		}
		c.logger.Infof("STATUS is: %v", c.cluster.Status.Members.Ready)
		k8sutil.SetEtcdVersion(pod, c.cluster.Spec.Version)
		patchdata, err := k8sutil.CreatePatch(oldpod, pod, v1.Pod{})
		if err != nil {
			return fmt.Errorf("error creating patch: %v", err)
		}
		_, err = c.config.KubeCli.CoreV1().Pods(ns).Patch(context.Background(), pod.GetName(), types.StrategicMergePatchType, patchdata, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("fail to update the etcd member (%s): %v", memberName, err)
		}
		c.logger.Infof("finished upgrading the etcd member %v", memberName)
		_, err = c.eventsCli.Create(context.Background(), k8sutil.MemberUpgradedEvent(memberName, k8sutil.GetEtcdVersion(oldpod), c.cluster.Spec.Version, c.cluster), metav1.CreateOptions{})
		if err != nil {
			c.logger.Errorf("failed to create member upgraded event: %v", err)
		}
		return nil
	}
	return nil
}

func (c *Cluster) rolloutOneMember(memberName, namespace string) error {
	c.logger.Infof("create a new member")
	err := c.addOneMember()
	if err != nil {
		c.logger.Warningf("unable add a new member")
		return err
	}
	//time.Sleep(60 * time.Second)
	listmembers, err := etcdutil.ListMembers(c.members.ClientURLs(), c.tlsConfig)
	if err != nil {
		c.logger.Warningf("unable list members")
		return err
	}
	var memberID uint64
	memberID = 0
	for _, member := range listmembers.Members {
		if member.Name == memberName {
			memberID = member.ID
		}
	}
	// excluir o pod
	err = c.removeMember(&etcdutil.Member{
		Name:      memberName,
		Namespace: namespace,
		ID:        memberID},
	)
	if err != nil {
		c.logger.Warningf("unable remove members")
		return err
	}
	return nil
}
