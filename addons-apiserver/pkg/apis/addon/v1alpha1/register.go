// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "addon.tanzu.vmware.com"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AntreaAddonConfig{},
		&AntreaAddonConfigList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

//var AddToScheme = func(scheme *runtime.Scheme) error {
//	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
//		Group:   "addon.tanzu.vmware.com",
//		Version: "v1alpha1",
//	})
//	// +kubebuilder:scaffold:install
//	scheme.AddKnownTypes(schema.GroupVersion{
//		Group:   "addon.tanzu.vmware.com",
//		Version: "v1alpha1",
//	}, &AntreaAddonConfig{}, &AntreaAddonConfigList{})
//	return nil
//}
