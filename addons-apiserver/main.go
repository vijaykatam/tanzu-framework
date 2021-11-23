// Copyright 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/klogr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apiserver"
	// TODO remove
	_ "k8s.io/code-generator"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	//klog.InitFlags(nil)

	// +kubebuilder:scaffold:scheme
}

func main() {

	ctrl.SetLogger(klogr.New())
	//setupLog.Info("Version", "version", buildinfo.Version, "buildDate", buildinfo.Date, "sha", buildinfo.SHA)

	setupLog.Info("start apiserver")
	a := apiserver.NewApiServer(ctrl.GetConfigOrDie(), setupLog)
	if err := a.Start(); err != nil {
		panic(err)
		os.Exit(1)
	}

	setupLog.Info("starting manager")

}
