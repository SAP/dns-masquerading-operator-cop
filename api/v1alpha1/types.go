/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and dns-masquerading-operator-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"

	"github.com/sap/component-operator-runtime/pkg/component"
	componentoperatorruntimetypes "github.com/sap/component-operator-runtime/pkg/types"
)

// DNSMasqueradingOperatorSpec defines the desired state of DNSMasqueradingOperator.
type DNSMasqueradingOperatorSpec struct {
	component.Spec `json:",inline"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	ReplicaCount int `json:"replicaCount,omitempty"`
	// +optional
	Image                          component.ImageSpec `json:"image"`
	component.KubernetesProperties `json:",inline"`
	EnableServiceController        bool `json:"enableServiceController,omitempty"`
	EnableIngressController        bool `json:"enableIngressController,omitempty"`
	EnableIstioGatewayController   bool `json:"enableIstioGatewayController,omitempty"`
}

// DNSMasqueradingOperatorStatus defines the observed state of DNSMasqueradingOperator.
type DNSMasqueradingOperatorStatus struct {
	component.Status `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +genclient

// DNSMasqueradingOperator is the Schema for the dnsmasqueradingoperators API.
type DNSMasqueradingOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DNSMasqueradingOperatorSpec `json:"spec,omitempty"`
	// +kubebuilder:default={"observedGeneration":-1}
	Status DNSMasqueradingOperatorStatus `json:"status,omitempty"`
}

var _ component.Component = &DNSMasqueradingOperator{}

// +kubebuilder:object:root=true

// DNSMasqueradingOperatorList contains a list of DNSMasqueradingOperator.
type DNSMasqueradingOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSMasqueradingOperator `json:"items"`
}

func (s *DNSMasqueradingOperatorSpec) ToUnstructured() map[string]any {
	result, err := runtime.DefaultUnstructuredConverter.ToUnstructured(s)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *DNSMasqueradingOperator) GetDeploymentNamespace() string {
	if c.Spec.Namespace != "" {
		return c.Spec.Namespace
	}
	return c.Namespace
}

func (c *DNSMasqueradingOperator) GetDeploymentName() string {
	if c.Spec.Name != "" {
		return c.Spec.Name
	}
	return c.Name
}

func (c *DNSMasqueradingOperator) GetSpec() componentoperatorruntimetypes.Unstructurable {
	return &c.Spec
}

func (c *DNSMasqueradingOperator) GetStatus() *component.Status {
	return &c.Status.Status
}

func init() {
	SchemeBuilder.Register(&DNSMasqueradingOperator{}, &DNSMasqueradingOperatorList{})
}
