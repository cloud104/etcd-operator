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

package k8sutil

import (
	"testing"

	api "github.com/cloud104/etcd-operator/pkg/apis/etcd/v1beta2"
	"github.com/cloud104/etcd-operator/pkg/util/etcdutil"
)

func TestDefaultBusyboxImageName(t *testing.T) {
	policy := &api.PodPolicy{}
	image := imageNameBusybox(policy)
	expected := defaultBusyboxImage
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestDefaultNilBusyboxImageName(t *testing.T) {
	image := imageNameBusybox(nil)
	expected := defaultBusyboxImage
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestSetBusyboxImageName(t *testing.T) {
	policy := &api.PodPolicy{
		BusyboxImage: "myRepo/busybox:1.3.2",
	}
	image := imageNameBusybox(policy)
	expected := "myRepo/busybox:1.3.2"
	if image != expected {
		t.Errorf("expect image=%s, get=%s", expected, image)
	}
}

func TestNewEtcdPodProbes(t *testing.T) {
	m := &etcdutil.Member{
		Name:      "k8s-test",
		Namespace: "tstetcd03",
	}
	initialCluster := []string{"k8s-test=http://k8s-test.k8s.tstetcd03.svc:2380"}

	tests := []struct {
		name              string
		version           string
		expectExecProbes  bool
		expectLivePath    string
		expectReadyPath   string
		expectStartupPath string
	}{
		{
			name:             "etcd < 3.5.7 uses exec probes without startupProbe",
			version:          "3.5.6",
			expectExecProbes: true,
		},
		{
			name:              "etcd 3.5.7..3.5.11 uses /health HTTP probes",
			version:           "3.5.9",
			expectLivePath:    "/health?exclude=NOSPACE&serializable=true",
			expectReadyPath:   "/health?serializable=false",
			expectStartupPath: "/health?exclude=NOSPACE&serializable=true",
		},
		{
			name:              "etcd >= 3.5.12 uses /livez and /readyz HTTP probes",
			version:           "3.5.16",
			expectLivePath:    "/livez",
			expectReadyPath:   "/readyz",
			expectStartupPath: "/livez",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cs := api.ClusterSpec{
				Version:    tc.version,
				Repository: "quay.io/coreos/etcd",
			}
			pod := newEtcdPod(m, initialCluster, "k8s", "new", "token", cs)
			c := pod.Spec.Containers[0]

			if tc.expectExecProbes {
				if c.LivenessProbe == nil || c.LivenessProbe.Exec == nil {
					t.Fatalf("expected exec livenessProbe for version %s", tc.version)
				}
				if c.ReadinessProbe == nil || c.ReadinessProbe.Exec == nil {
					t.Fatalf("expected exec readinessProbe for version %s", tc.version)
				}
				if c.StartupProbe != nil {
					t.Fatalf("expected nil startupProbe for version %s", tc.version)
				}
				return
			}

			if c.LivenessProbe == nil || c.LivenessProbe.HTTPGet == nil || c.LivenessProbe.HTTPGet.Path != tc.expectLivePath {
				t.Errorf("expected livenessProbe path=%s, got %+v", tc.expectLivePath, c.LivenessProbe)
			}
			if c.ReadinessProbe == nil || c.ReadinessProbe.HTTPGet == nil || c.ReadinessProbe.HTTPGet.Path != tc.expectReadyPath {
				t.Errorf("expected readinessProbe path=%s, got %+v", tc.expectReadyPath, c.ReadinessProbe)
			}
			if c.StartupProbe == nil || c.StartupProbe.HTTPGet == nil || c.StartupProbe.HTTPGet.Path != tc.expectStartupPath {
				t.Fatalf("expected startupProbe path=%s, got %+v", tc.expectStartupPath, c.StartupProbe)
			}
			if c.StartupProbe.FailureThreshold != 24 || c.StartupProbe.PeriodSeconds != 10 {
				t.Errorf("expected startupProbe failureThreshold=24 and periodSeconds=10, got %d and %d",
					c.StartupProbe.FailureThreshold, c.StartupProbe.PeriodSeconds)
			}
		})
	}
}
