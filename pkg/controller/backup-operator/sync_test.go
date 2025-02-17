// Copyright 2018 The etcd-operator Authors
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
	"testing"

	api "github.com/cloud104/etcd-operator/pkg/apis/etcd/v1beta2"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		spec      *api.BackupSpec
		expectErr bool
	}{{
		spec: &api.BackupSpec{
			EtcdEndpoints: []string{"http://localhost:2379"},
		},
		expectErr: false,
	}, { // fail due to empty etcd endpoints
		spec:      &api.BackupSpec{},
		expectErr: true,
	}}

	for i, tt := range tests {
		err := validate(tt.spec)
		if err != nil && !tt.expectErr {
			t.Errorf("#%d: validate failed: %v", i, err)
		}
		if err == nil && tt.expectErr {
			t.Errorf("#%d: expect error, but got nil", i)
		}
	}
}
