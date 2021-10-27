#!/usr/bin/env bash

# Copyright 2021 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

go mod vendor
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
#CODEGEN_PKG=/Users/vkatam/go/pkg/mod/k8s.io/code-generator@v0.20.5
# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}/generate-groups.sh" all \
  github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/generated github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis \
  "addon:v1alpha1" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

bash "${CODEGEN_PKG}/generate-internal-groups.sh" "deepcopy,defaulter,conversion,openapi" \
  github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/generated github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis \
  "addon:v1alpha1" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
rm -rf vendor/