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

package k8sutil

import (
	"context"
	"encoding/json"
	"fmt"
	api "github.com/cloud104/etcd-operator/pkg/apis/etcd/v1beta2"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

// TODO: replace this package with Operator client

// EtcdClusterCRUpdateFunc is a function to be used when atomically
// updating a Cluster CR.
type EtcdClusterCRUpdateFunc func(*api.EtcdCluster)

func GetClusterList(restcli rest.Interface, ns string) (*api.EtcdClusterList, error) {
	b, err := restcli.Get().RequestURI(listClustersURI(ns)).DoRaw(context.Background())
	if err != nil {
		return nil, err
	}

	clusters := &api.EtcdClusterList{}
	if err := json.Unmarshal(b, clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func listClustersURI(ns string) string {
	return fmt.Sprintf("/apis/%s/namespaces/%s/%s", api.SchemeGroupVersion.String(), ns, api.EtcdClusterResourcePlural)
}

//func CreateCRD(clientset apiextensionsclient.Interface, crdName, rkind, rplural, shortName string) error {
//	crd := &apiextensionsv1.CustomResourceDefinition{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: crdName,
//		},
//		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
//			Group: api.SchemeGroupVersion.Group,
//			Names: apiextensionsv1.CustomResourceDefinitionNames{
//				Plural: rplural,
//				Kind:   rkind,
//			},
//			Scope: apiextensionsv1.NamespaceScoped,
//			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
//				{
//					Name:       "v1beta2",
//					Served:     true,
//					Storage:    true,
//					Deprecated: false,
//					Schema: &apiextensionsv1.CustomResourceValidation{
//						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
//							Type: "object",
//							Properties: map[string]apiextensionsv1.JSONSchemaProps{
//								"spec": {
//									Type: "object",
//									Properties: map[string]apiextensionsv1.JSONSchemaProps{
//										"size": {
//											Type: "integer",
//										},
//										"version": {
//											Type: "string",
//										},
//										"repository": {
//											Type: "string",
//										},
//									},
//								},
//							},
//						},
//					},
//					Subresources: &apiextensionsv1.CustomResourceSubresources{
//						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
//					},
//				},
//			},
//			PreserveUnknownFields: false,
//		},
//	}
//	if len(shortName) != 0 {
//		crd.Spec.Names.ShortNames = []string{shortName}
//	}
//	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
//	if err != nil && !IsKubernetesResourceAlreadyExistError(err) {
//		return err
//	}
//	return nil
//}

//func WaitCRDReady(clientset apiextensionsclient.Interface, crdName string) error {
//	err := retryutil.Retry(5*time.Second, 20, func() (bool, error) {
//		crd, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), crdName, metav1.GetOptions{})
//		if err != nil {
//			return false, err
//		}
//		for _, cond := range crd.Status.Conditions {
//			switch cond.Type {
//			case apiextensionsv1.Established:
//				if cond.Status == apiextensionsv1.ConditionTrue {
//					return true, nil
//				}
//			case apiextensionsv1.NamesAccepted:
//				if cond.Status == apiextensionsv1.ConditionFalse {
//					return false, fmt.Errorf("Name conflict: %v", cond.Reason)
//				}
//			}
//		}
//		return false, nil
//	})
//	if err != nil {
//		return fmt.Errorf("wait CRD created failed: %v", err)
//	}
//	return nil
//}

func MustNewKubeExtClient() apiextensionsclient.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
	return apiextensionsclient.NewForConfigOrDie(cfg)
}
