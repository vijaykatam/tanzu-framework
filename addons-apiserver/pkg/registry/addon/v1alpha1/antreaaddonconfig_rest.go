package v1alpha1

import (
	"context"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis/addon/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
)

type AntreaAddonConfigREST struct {
	client kubernetes.Interface
}

var (
	_ rest.StandardStorage    = &AntreaAddonConfigREST{}
	_ rest.ShortNamesProvider = &AntreaAddonConfigREST{}
)

func NewAntreaAddonConfigREST(client kubernetes.Interface) *AntreaAddonConfigREST {
	return &AntreaAddonConfigREST{client: client}
}

func (r *AntreaAddonConfigREST) ShortNames() []string {
	return []string{"aac"}
}

func (r *AntreaAddonConfigREST) NamespaceScoped() bool {
	return true
}

func (r *AntreaAddonConfigREST) New() runtime.Object {
	return &v1alpha1.AntreaAddonConfig{}
}

func (r *AntreaAddonConfigREST) NewList() runtime.Object {
	return &v1alpha1.AntreaAddonConfigList{}
}

func (r *AntreaAddonConfigREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return &v1alpha1.AntreaAddonConfig{}, nil
}

func (r *AntreaAddonConfigREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return &v1alpha1.AntreaAddonConfig{}, nil
}
func (r *AntreaAddonConfigREST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return &v1alpha1.AntreaAddonConfigList{}, nil
}
func (r *AntreaAddonConfigREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return &v1alpha1.AntreaAddonConfig{}, false, nil
}

func (r *AntreaAddonConfigREST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	return &v1alpha1.AntreaAddonConfig{}, false, nil
}

func (r *AntreaAddonConfigREST) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	return &v1alpha1.AntreaAddonConfigList{}, nil
}

func (r *AntreaAddonConfigREST) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rest.NewDefaultTableConvertor(v1alpha1.Resource("antreaaddonconfigs")).
		ConvertToTable(ctx, obj, tableOptions)
}

func (r *AntreaAddonConfigREST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return nil, nil
}
