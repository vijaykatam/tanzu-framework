// Copyright 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apiserver"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/klogr"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
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
		setupLog.Error(err, "error starting apiserver")
		os.Exit(1)
	}

	setupLog.Info("starting manager")

}
