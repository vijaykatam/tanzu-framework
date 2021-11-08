package generators

import (
	"fmt"
	"io"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	"path/filepath"
	"reflect"
	"strings"
)

// CustomArgs is used tby the go2idl framework to pass args specific to this
// generator.
type CustomArgs struct {
	BoundingDirs []string // Only deal with types rooted under these dirs.
}

// These are the comment tags that carry parameters for datavalue conversion.
const (
	tagEnabledName   = "tanzu-framework:data-value-converter-gen"
	dataValueTagName = "dataValue"
)

const tagValuePackage = "package"

type enabledTagValue struct {
	value    string
	register bool
}

func extractDataValueTag(comments []string) []string {
	return types.ExtractCommentTags("+", comments)[dataValueTagName]
}

func extractEnabledTypeTag(t *types.Type) *enabledTagValue {
	comments := append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
	return extractEnabledTag(comments)
}

func extractEnabledTag(comments []string) *enabledTagValue {
	tagVals := types.ExtractCommentTags("+", comments)[tagEnabledName]
	if tagVals == nil {
		// No match for the tag.
		return nil
	}
	// If there are multiple values, abort.
	if len(tagVals) > 1 {
		klog.Fatalf("Found %d %s tags: %q", len(tagVals), tagEnabledName, tagVals)
	}

	// If we got here we are returning something.
	tag := &enabledTagValue{}

	// Get the primary value.
	parts := strings.Split(tagVals[0], ",")
	if len(parts) >= 1 {
		tag.value = parts[0]
	}

	// Parse extra arguments.
	parts = parts[1:]
	for i := range parts {
		kv := strings.SplitN(parts[i], "=", 2)
		k := kv[0]
		v := ""
		if len(kv) == 2 {
			v = kv[1]
		}
		switch k {
		case "register":
			if v != "false" {
				tag.register = true
			}
		default:
			klog.Fatalf("Unsupported %s param: %q", tagEnabledName, parts[i])
		}
	}
	return tag
}

func checkTag(comments []string, require ...string) bool {
	values := types.ExtractCommentTags("+", comments)[tagEnabledName]
	if len(require) == 0 {
		return len(values) == 1 && values[0] == ""
	}
	return reflect.DeepEqual(values, require)
}

func fromDataValue() *namer.NameStrategy {
	return &namer.NameStrategy{
		Prefix: "FromDataValue_",
		Join: func(pre string, in []string, post string) string {
			return pre + strings.Join(in, "_") + post
		},
	}
}

func toDataValue() *namer.NameStrategy {
	return &namer.NameStrategy{
		Prefix: "ToDataValue_",
		Join: func(pre string, in []string, post string) string {
			return pre + strings.Join(in, "_") + post
		},
	}
}

func objectDefaultFnNamer() *namer.NameStrategy {
	return &namer.NameStrategy{
		Prefix: "SetObjectDefaults_",
		Join: func(pre string, in []string, post string) string {
			return pre + strings.Join(in, "_") + post
		},
	}
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":          namer.NewPublicNamer(1),
		"private":         namer.NewPrivateNamer(0),
		"raw":             namer.NewRawNamer("", nil),
		"fromdatavaluefn": fromDataValue(),
		"todatavaluefn":   toDataValue(),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)), boilerplate...)

	boundingDirs := []string{}
	if customArgs, ok := arguments.CustomArgs.(*CustomArgs); ok {
		if customArgs.BoundingDirs == nil {
			customArgs.BoundingDirs = context.Inputs
		}
		for i := range customArgs.BoundingDirs {
			// Strip any trailing slashes - they are not exactly "correct" but
			// this is friendlier.
			boundingDirs = append(boundingDirs, strings.TrimRight(customArgs.BoundingDirs[i], "/"))
		}
	}

	for i := range inputs {
		klog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		ptag := extractEnabledTag(pkg.Comments)
		ptagValue := ""
		ptagRegister := false
		if ptag != nil {
			ptagValue = ptag.value
			if ptagValue != tagValuePackage {
				klog.Fatalf("Package %v: unsupported %s value: %q", i, tagEnabledName, ptagValue)
			}
			ptagRegister = ptag.register
			klog.V(5).Infof("  tag.value: %q, tag.register: %t", ptagValue, ptagRegister)
		} else {
			klog.V(5).Infof("  no tag")
		}

		// If the pkg-scoped tag says to generate, we can skip scanning types.
		pkgNeedsGeneration := ptagValue == tagValuePackage
		if !pkgNeedsGeneration {
			// If the pkg-scoped tag did not exist, scan all types for one that
			// explicitly wants generation.
			for _, t := range pkg.Types {
				klog.V(5).Infof("  considering type %q", t.Name.String())
				ttag := extractEnabledTypeTag(t)
				if ttag != nil && ttag.value == "true" {
					klog.V(5).Infof("    tag=true")
					pkgNeedsGeneration = true
					break
				}
			}
		}

		if pkgNeedsGeneration {
			klog.V(3).Infof("Package %q needs generation", i)
			path := pkg.Path
			// if the source path is within a /vendor/ directory (for example,
			// k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/apis/meta/v1), allow
			// generation to output to the proper relative path (under vendor).
			// Otherwise, the generator will create the file in the wrong location
			// in the output directory.
			// TODO: build a more fundamental concept in gengo for dealing with modifications
			// to vendored packages.
			if strings.HasPrefix(pkg.SourcePath, arguments.OutputBase) {
				expandedPath := strings.TrimPrefix(pkg.SourcePath, arguments.OutputBase)
				if strings.Contains(expandedPath, "/vendor/") {
					path = expandedPath
				}
			}
			packages = append(packages,
				&generator.DefaultPackage{
					PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
					PackagePath: path,
					HeaderText:  header,
					GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
						return []generator.Generator{
							NewGenDataValueConverter(arguments.OutputFileBaseName, pkg.Path, boundingDirs, ptagValue == tagValuePackage, ptagRegister),
						}
					},
					FilterFunc: func(c *generator.Context, t *types.Type) bool {
						return t.Name.Package == pkg.Path
					},
				})
		}
	}
	return packages
}

// genSet produces a file with a set for a single type.
type genDataValueConverter struct {
	generator.DefaultGen
	targetPackage string
	boundingDirs  []string
	allTypes      bool
	registerTypes bool
	imports       namer.ImportTracker
	typesForInit  []*types.Type
}

func NewGenDataValueConverter(sanitizedName, targetPackage string, boundingDirs []string, allTypes, registerTypes bool) generator.Generator {
	return &genDataValueConverter{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		targetPackage: targetPackage,
		boundingDirs:  boundingDirs,
		allTypes:      allTypes,
		registerTypes: registerTypes,
		imports:       generator.NewImportTracker(),
		typesForInit:  make([]*types.Type, 0),
	}
}

func (g *genDataValueConverter) Filter(c *generator.Context, t *types.Type) bool {
	// Filter out types not being processed or not copyable within the package.
	enabled := g.allTypes
	if !enabled {
		ttag := extractEnabledTypeTag(t)
		if ttag != nil && ttag.value == "true" {
			enabled = true
		}
	}
	if !enabled {
		return false
	}

	g.typesForInit = append(g.typesForInit, t)
	return true
}

func isRootedUnder(pkg string, roots []string) bool {
	// Add trailing / to avoid false matches, e.g. foo/bar vs foo/barn.  This
	// assumes that bounding dirs do not have trailing slashes.
	pkg = pkg + "/"
	for _, root := range roots {
		if strings.HasPrefix(pkg, root+"/") {
			return true
		}
	}
	return false
}

// deepCopyMethod returns the signature of a DeepCopy() method, nil or an error
// if the type does not match. This allows more efficient deep copy
// implementations to be defined by the type's author.  The correct signature
//    func (t *T) FromDataValue(data []byte)
func fromDataValueMethod(t *types.Type) (*types.Signature, error) {
	f, found := t.Methods["FromDataValue"]
	if !found {
		return nil, nil
	}
	if len(f.Signature.Parameters) != 1 {
		return nil, fmt.Errorf("type %v: invalid FromDataValue signature, expected 1 parameter", t)
	}
	if len(f.Signature.Results) != 0 {
		return nil, fmt.Errorf("type %v: invalid FromDataValue signature, expected no result", t)
	}

	ptrRcvr := f.Signature.Receiver != nil && f.Signature.Receiver.Kind == types.Pointer && f.Signature.Receiver.Elem.Name == t.Name

	if !ptrRcvr {
		return nil, fmt.Errorf("type %v: invalid FromDataValue signature, expected a *%s receiver", t, t.Name.Name)
	}

	return f.Signature, nil
}


// fromDataValueMethodOrDie returns the signatrue of a DeepCopy method, nil or calls klog.Fatalf
// if the type does not match.
func fromDataValueMethodOrDie(t *types.Type) *types.Signature {
	ret, err := fromDataValueMethod(t)
	if err != nil {
		klog.Fatal(err)
	}
	return ret
}

func underlyingType(t *types.Type) *types.Type {
	for t.Kind == types.Alias {
		t = t.Underlying
	}
	return t
}

func (g *genDataValueConverter) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genDataValueConverter) Imports(c *generator.Context) (imports []string) {
	importLines := []string{}
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return append(importLines, "github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/yttdatavalues")
}

func argsFromType(ts ...*types.Type) generator.Args {
	a := generator.Args{
		"type": ts[0],
	}
	for i, t := range ts {
		a[fmt.Sprintf("type%d", i+1)] = t
	}
	return a
}

func (g *genDataValueConverter) Init(c *generator.Context, w io.Writer) error {
	return nil
}

func (g *genDataValueConverter) needsGeneration(t *types.Type) bool {
	tag := extractEnabledTypeTag(t)
	tv := ""
	if tag != nil {
		tv = tag.value
		if tv != "true" && tv != "false" {
			klog.Fatalf("Type %v: unsupported %s value: %q", t, tagEnabledName, tag.value)
		}
	}
	if g.allTypes && tv == "false" {
		// The whole package is being generated, but this type has opted out.
		klog.V(5).Infof("Not generating for type %v because type opted out", t)
		return false
	}
	if !g.allTypes && tv != "true" {
		// The whole package is NOT being generated, and this type has NOT opted in.
		klog.V(5).Infof("Not generating for type %v because type did not opt in", t)
		return false
	}
	return true
}
