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
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AntreaAddonConfigSpec `json:"spec"`
	// +optional
	Status AntreaAddonConfigStatus `json:"status,omitempty"`
}

// AntreaAddonConfigList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AntreaAddonConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AntreaAddonConfig `json:"items"`
}

// AntreaAddonConfigSpec defines the desired state of AntreaAddonConfig
// +tanzu-framework:data-value-converter-gen=true
type AntreaAddonConfigSpec struct {
	InfraProvider string `json:"infraProvider" dataValue:"infraProvider"`

	// +default="100.64.0.0/13"
	// +dataValue="antrea.config.serviceCIDR"
	ServiceCIDR *string `json:"serviceCIDR"`

	// +default="encap"
	// +dataValue="antrea.config.trafficEncapMode"
	TrafficEncapMode *string `json:"trafficEncapMode"`

	// +default=false
	// +dataValue="antrea.config.noSNAT"
	NoSNAT *bool `json:"noSNAT"`

	// +dataValue="antrea.config.defaultMTU"
	DefaultMTU string `json:"defaultMTU,omitempty"`

	// +default="TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384"
	// +dataValue="antrea.config.tlsCipherSuites"
	TlsCipherSuites *string `json:"tlsCipherSuites"`

	// +dataValue="antrea.config.featureGates"
	FeatureGates *FeatureGates `json:"featureGates"`
}

type FeatureGates struct {
	// +default=true
	// +dataValue="antrea.config.featureGates.AntreaProxy"
	AntreaProxy *bool `json:"antreaProxy"`

	// +default=false
	// +dataValue="antrea.config.featureGates.EndpointSlice"
	EndpointSlice *bool `json:"endpointSlice"`

	// +default=true
	// +dataValue="antrea.config.featureGates.AntreaPolicy"
	AntreaPolicy *bool `json:"antreaPolicy"`

	// +default=false
	// +dataValue="antrea.config.featureGates.NodePortLocal"
	NodePortLocal *bool `json:"nodePortLocal"`

	// +default=true
	// +dataValue="antrea.config.featureGates.AntreaTraceflow"
	AntreaTraceflow *bool `json:"antreaTraceflow"`

	// +default=false
	// +dataValue="antrea.config.featureGates.FlowExporter"
	FlowExporter *bool `json:"flowExporter`

	// +default=false
	// +dataValue="antrea.config.featureGates.NetworkPolicyStats"
	NetworkPolicyStats *bool `json:"networkPolicyStats"`
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
