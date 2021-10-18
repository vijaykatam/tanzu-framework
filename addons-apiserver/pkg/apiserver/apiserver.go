package apiserver

import (
	"context"
	"github.com/go-logr/logr"
	addonconfigv1alpha1 "github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/apis/addon/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aggregatorclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"net"
	"os"
	"path/filepath"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"
	"github.com/vmware-tanzu/tanzu-framework/addons-apiserver/pkg/generated/openapi"
)

type apiServer struct {
	config *rest.Config
	logger logr.Logger
}

func NewApiServer(config *rest.Config, logger logr.Logger) *apiServer {
	return &apiServer{config: config, logger: logger}
}

func (a *apiServer) Start() error {
	server := builder.APIServer.
		WithResource(&addonconfigv1alpha1.AntreaAddonConfig{}).
		WithOpenAPIDefinitions("addons-apiserver", "v1alpha1", openapi.GetOpenAPIDefinitions).
		WithoutEtcd()
	//server, err := withAPIServiceAndSelfSignedCerts(server)
	//if err != nil {
	//	return err
	//}
	//server.WithConfigFns()
	server.WithOptionsFns(a.selfSignedCertsServerOption)
	//
	//server.WithServerFns()
	//cmd, err := server.Build()
	//builder.ServerOptions{}

	// TODO: cert rotate

	return server.Execute()

}

func (a *apiServer) selfSignedCertsServerOption(o *builder.ServerOptions) *builder.ServerOptions {
	o.RecommendedOptions.SecureServing.ServerCert.CertDirectory = filepath.Join(homedir.HomeDir(), "addons-apiserver")
	o.RecommendedOptions.SecureServing.ServerCert.PairName = "addons-apiserver"
	o.RecommendedOptions.SecureServing.BindPort = 10299
	if err := o.RecommendedOptions.SecureServing.
		MaybeDefaultWithSelfSignedCerts(
			"addons-apiserver",
			[]string{"addons-apiserver.tkg-system.svc"},
			[]net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		a.logger.Error(err, "error creating cert")
		os.Exit(1)
	}
	aggClient, err := aggregatorclient.NewForConfig(a.config)
	if err != nil {
		a.logger.Error(err, "error creating cert")
		os.Exit(1)
	}
	caContentProvider, err := dynamiccertificates.NewDynamicCAContentFromFile(
		"self-signed cert",
		o.RecommendedOptions.SecureServing.ServerCert.CertKey.CertFile)
	if err != nil {
		a.logger.Error(err, "error retrieving cert")
		os.Exit(1)
	}
	apiService, err := aggClient.ApiregistrationV1().APIServices().Get(context.TODO(), "v1alpha1.addon.tanzu.vmware.com", metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			addonsService := &apiregv1.APIService{
				ObjectMeta: metav1.ObjectMeta{
					Name: "v1alpha1.addon.tanzu.vmware.com",
				},
				Spec: apiregv1.APIServiceSpec{
					Group:                "addon.tanzu.vmware.com",
					GroupPriorityMinimum: 100,
					Version:              "v1alpha1",
					VersionPriority:      100,
					Service: &apiregv1.ServiceReference{
						Name:      "addons-apiserver",
						Namespace: "tkg-system",
					},
					CABundle: caContentProvider.CurrentCABundleContent(),
				},
			}
			_, err := aggClient.ApiregistrationV1().APIServices().Create(context.TODO(), addonsService, metav1.CreateOptions{})
			if err != nil {
				a.logger.Error(err, "error creating apiservice")
				os.Exit(1)
			}
		} else {
			a.logger.Error(err, "error retrieving apiservice")
			os.Exit(1)
		}
	} else {
		apiService.Spec.CABundle = caContentProvider.CurrentCABundleContent()
		_, err := aggClient.ApiregistrationV1().APIServices().Update(context.TODO(), apiService, metav1.UpdateOptions{})
		if err != nil {
			a.logger.Error(err, "error updating apiservice")
			os.Exit(1)
		}
	}
	return o
}
