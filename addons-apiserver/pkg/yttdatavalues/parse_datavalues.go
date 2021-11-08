package yttdatavalues

import (
	"github.com/jeremywohl/flatten"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
	"github.com/k14s/ytt/pkg/workspace"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func GetFlattenedDataValues(data []byte) (map[string]interface{}, error) {
	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", data)),
	})
	ui := ui.NewTTY(false)

	lib := workspace.NewRootLibrary(filesToProcess)
	libCtx := workspace.LibraryExecutionContext{Current: lib, Root: lib}

	libExecFact := workspace.NewLibraryExecutionFactory(&ui, workspace.TemplateLoaderOpts{})
	rootLibraryExecution := libExecFact.New(libCtx)
	schema, _, err := rootLibraryExecution.Schemas(nil)
	if err != nil {
		return nil, err
	}

	values, _, err := rootLibraryExecution.Values([]*workspace.DataValues{}, schema)
	if err != nil || values == nil || values.Doc == nil {
		return nil, errors.Wrap(err, "unable to load yaml document")
	}
	yamlBytes, err := values.Doc.AsYAMLBytes()
	if err != nil {
		return nil, err
	}

	dataValues := make(map[string]interface{})
	if err := yaml.Unmarshal(yamlBytes, &dataValues); err != nil {
		return nil, err
	}

	return flatten.Flatten(dataValues, "", flatten.DotStyle)
}
