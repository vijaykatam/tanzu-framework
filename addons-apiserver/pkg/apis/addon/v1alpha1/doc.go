// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis/addon
// +k8s:defaulter-gen=TypeMeta
// +groupName=addon.tanzu.vmware.com
package v1alpha1 // import "github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis/addon/v1alpha1"
