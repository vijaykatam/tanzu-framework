// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AntreaAddonConfig
// +k8s:openapi-gen=true
type AntreaAddonConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AntreaAddonConfigSpec   `json:"spec,omitempty"`
	Status AntreaAddonConfigStatus `json:"status,omitempty"`
}

// AntreaAddonConfigList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AntreaAddonConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AntreaAddonConfig `json:"items"`
}

// AntreaAddonConfigSpec defines the desired state of AntreaAddonConfig
type AntreaAddonConfigSpec struct {
	InfraProvider string `json:"infraProvider"`
	ServiceCIDR   string `json:"serviceCIDR"`
	ServiceCIDRv6 string `json:"serviceCIDRv6"`
}

var _ resource.Object = &AntreaAddonConfig{}
var _ resourcestrategy.Validater = &AntreaAddonConfig{}

func (in *AntreaAddonConfig) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *AntreaAddonConfig) NamespaceScoped() bool {
	return false
}

func (in *AntreaAddonConfig) New() runtime.Object {
	return &AntreaAddonConfig{}
}

func (in *AntreaAddonConfig) NewList() runtime.Object {
	return &AntreaAddonConfigList{}
}

func (in *AntreaAddonConfig) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "addon.tanzu.vmware.com",
		Version:  "v1alpha1",
		Resource: "antreaaddonconfigs",
	}
}

func (in *AntreaAddonConfig) IsStorageVersion() bool {
	return true
}

func (in *AntreaAddonConfig) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &AntreaAddonConfigList{}

func (in *AntreaAddonConfigList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// AntreaAddonConfigStatus defines the observed state of AntreaAddonConfig
type AntreaAddonConfigStatus struct {
}

func (in AntreaAddonConfigStatus) SubResourceName() string {
	return "status"
}

// AntreaAddonConfig implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &AntreaAddonConfig{}

func (in *AntreaAddonConfig) GetStatus() resource.StatusSubResource {
	return in.Status
}

// AntreaAddonConfigStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &AntreaAddonConfigStatus{}

func (in AntreaAddonConfigStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*AntreaAddonConfig).Status = in
}
