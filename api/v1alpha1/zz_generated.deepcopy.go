//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and dns-masquerading-operator-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSMasqueradingOperator) DeepCopyInto(out *DNSMasqueradingOperator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSMasqueradingOperator.
func (in *DNSMasqueradingOperator) DeepCopy() *DNSMasqueradingOperator {
	if in == nil {
		return nil
	}
	out := new(DNSMasqueradingOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSMasqueradingOperator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSMasqueradingOperatorList) DeepCopyInto(out *DNSMasqueradingOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DNSMasqueradingOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSMasqueradingOperatorList.
func (in *DNSMasqueradingOperatorList) DeepCopy() *DNSMasqueradingOperatorList {
	if in == nil {
		return nil
	}
	out := new(DNSMasqueradingOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSMasqueradingOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSMasqueradingOperatorSpec) DeepCopyInto(out *DNSMasqueradingOperatorSpec) {
	*out = *in
	out.Spec = in.Spec
	out.Image = in.Image
	in.KubernetesProperties.DeepCopyInto(&out.KubernetesProperties)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSMasqueradingOperatorSpec.
func (in *DNSMasqueradingOperatorSpec) DeepCopy() *DNSMasqueradingOperatorSpec {
	if in == nil {
		return nil
	}
	out := new(DNSMasqueradingOperatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSMasqueradingOperatorStatus) DeepCopyInto(out *DNSMasqueradingOperatorStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSMasqueradingOperatorStatus.
func (in *DNSMasqueradingOperatorStatus) DeepCopy() *DNSMasqueradingOperatorStatus {
	if in == nil {
		return nil
	}
	out := new(DNSMasqueradingOperatorStatus)
	in.DeepCopyInto(out)
	return out
}
