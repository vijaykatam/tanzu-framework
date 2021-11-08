// +build tools

// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import _ "k8s.io/code-generator"
