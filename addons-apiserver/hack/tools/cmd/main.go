package main

import (
	"fmt"
	"os"

	"github.com/jeremywohl/flatten"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
	"github.com/k14s/ytt/pkg/workspace"
	"gopkg.in/yaml.v3"
)

const antreaDataValue = `
#@data/values
#@overlay/match-child-defaults missing_ok=True
---
infraProvider: vsphere
antrea:
  config:
    serviceCIDR: 100.64.0.0/13
    trafficEncapMode: encap
    noSNAT: false
    tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
    featureGates:
      AntreaProxy: true
      EndpointSlice: false
      AntreaPolicy: true
      NodePortLocal: false
      AntreaTraceflow: true
`

func main() {
	yttDataValues := []byte(antreaDataValue)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yttDataValues)),
	})
	ui := ui.NewTTY(false)

	lib := workspace.NewRootLibrary(filesToProcess)
	libCtx := workspace.LibraryExecutionContext{Current: lib, Root: lib}

	libExecFact := workspace.NewLibraryExecutionFactory(&ui, workspace.TemplateLoaderOpts{})
	rootLibraryExecution := libExecFact.New(libCtx)
	schema, _, err := rootLibraryExecution.Schemas(nil)
	if err != nil {
		fmt.Errorf("unable to load schemas %v", err)
		os.Exit(1)
	}

	values, _, err := rootLibraryExecution.Values([]*workspace.DataValues{}, schema)
	if err != nil || values == nil || values.Doc == nil {
		fmt.Errorf("unable to load yaml document %v", err)
		os.Exit(1)

	}
	yamlBytes, err := values.Doc.AsYAMLBytes()
	if err != nil {
		fmt.Errorf("unable to load yaml document %v", err)
		os.Exit(1)
	}

	dataValues := make(map[string]interface{})
	if err := yaml.Unmarshal(yamlBytes, &dataValues); err != nil {
		fmt.Errorf("unable to load yaml document %v", err)
		os.Exit(1)
	}

	dataValuesFlattened, err := flatten.Flatten(dataValues, "", flatten.DotStyle)
	if err != nil {
		fmt.Errorf("unable to load yaml document %v", err)
		os.Exit(1)
	}
	for k, v := range dataValuesFlattened {
		fmt.Println(fmt.Sprintf("%s:%v", k, v))
	}
}
