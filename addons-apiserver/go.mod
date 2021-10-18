module github.com/vmware-tanzu/tanzu-framework/addons-apiserver

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	k8s.io/apimachinery v0.22.2
	k8s.io/apiserver v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/code-generator v0.22.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.22.1
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e
	sigs.k8s.io/apiserver-runtime v1.0.3-0.20210913073608-0663f60bfee2
	sigs.k8s.io/controller-runtime v0.10.1
)
