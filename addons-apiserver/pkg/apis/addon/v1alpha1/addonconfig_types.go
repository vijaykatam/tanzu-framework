package v1alpha1

import (
	"context"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinitionSpec describes how a user wants their resource to appear
type AddonConfigSpec struct {
	// Group is the group this resource belongs in
	Group string

	// Names are the names used to describe this custom resource
	Names CustomResourceDefinitionNames

	// Validation describes the validation methods for CustomResources
	// Optional, the global validation schema for all versions.
	// Top-level and per-version schemas are mutually exclusive.
	// +optional
	Validation *CustomResourceValidation

	// Schema describes the schema for CustomResource used in validation, pruning, and defaulting.
	// Top-level and per-version schemas are mutually exclusive.
	// Per-version schemas must not all be set to identical values (top-level validation schema should be used instead)
	// This field is alpha-level and is only honored by servers that enable the CustomResourceWebhookConversion feature.
	// +optional
	Schema *CustomResourceValidation

	// AdditionalPrinterColumns are additional columns shown e.g. in kubectl next to the name. Defaults to a created-at column.
	// Optional, the global columns for all versions.
	// Top-level and per-version columns are mutually exclusive.
	// +optional
	AdditionalPrinterColumns []CustomResourceColumnDefinition

	// preserveUnknownFields disables pruning of object fields which are not
	// specified in the OpenAPI schema. apiVersion, kind, metadata and known
	// fields inside metadata are always preserved.
	// Defaults to true in v1beta and will default to false in v1.
	PreserveUnknownFields *bool
}

// CustomResourceColumnDefinition specifies a column for server side printing.
type CustomResourceColumnDefinition struct {
	// name is a human readable name for the column.
	Name string
	// type is an OpenAPI type definition for this column.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	Type string
	// format is an optional OpenAPI type definition for this column. The 'name' format is applied
	// to the primary identifier column to assist in clients identifying column is the resource name.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	Format string
	// description is a human readable description of this column.
	Description string
	// priority is an integer defining the relative importance of this column compared to others. Lower
	// numbers are considered higher priority. Columns that may be omitted in limited space scenarios
	// should be given a higher priority.
	Priority int32

	// JSONPath is a simple JSON path, i.e. without array notation.
	JSONPath string
}

// CustomResourceDefinitionNames indicates the names to serve this CustomResourceDefinition
type CustomResourceDefinitionNames struct {
	// Plural is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration
	// too: plural.group and it must be all lowercase.
	Plural string
	// Singular is the singular name of the resource.  It must be all lowercase  Defaults to lowercased <kind>
	Singular string
	// ShortNames are short names for the resource.  It must be all lowercase.
	ShortNames []string
	// Kind is the serialized kind of the resource.  It is normally CamelCase and singular.
	Kind string
	// ListKind is the serialized kind of the list for this resource.  Defaults to <kind>List.
	ListKind string
	// Categories is a list of grouped resources custom resources belong to (e.g. 'all')
	// +optional
	Categories []string
}

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// CustomResourceDefinitionConditionType is a valid value for CustomResourceDefinitionCondition.Type
type CustomResourceDefinitionConditionType string

const (
	// Established means that the resource has become active. A resource is established when all names are
	// accepted without a conflict for the first time. A resource stays established until deleted, even during
	// a later NamesAccepted due to changed names. Note that not all names can be changed.
	Established CustomResourceDefinitionConditionType = "Established"
	// NamesAccepted means the names chosen for this CustomResourceDefinition do not conflict with others in
	// the group and are therefore accepted.
	NamesAccepted CustomResourceDefinitionConditionType = "NamesAccepted"
	// NonStructuralSchema means that one or more OpenAPI schema is not structural.
	//
	// A schema is structural if it specifies types for all values, with the only exceptions of those with
	// - x-kubernetes-int-or-string: true — for fields which can be integer or string
	// - x-kubernetes-preserve-unknown-fields: true — for raw, unspecified JSON values
	// and there is no type, additionalProperties, default, nullable or x-kubernetes-* vendor extenions
	// specified under allOf, anyOf, oneOf or not.
	//
	// Non-structural schemas will not be allowed anymore in v1 API groups. Moreover, new features will not be
	// available for non-structural CRDs:
	// - pruning
	// - defaulting
	// - read-only
	// - OpenAPI publishing
	// - webhook conversion
	NonStructuralSchema CustomResourceDefinitionConditionType = "NonStructuralSchema"
	// Terminating means that the CustomResourceDefinition has been deleted and is cleaning up.
	Terminating CustomResourceDefinitionConditionType = "Terminating"
	// KubernetesAPIApprovalPolicyConformant indicates that an API in *.k8s.io or *.kubernetes.io is or is not approved.  For CRDs
	// outside those groups, this condition will not be set.  For CRDs inside those groups, the condition will
	// be true if .metadata.annotations["api-approved.kubernetes.io"] is set to a URL, otherwise it will be false.
	// See https://github.com/kubernetes/enhancements/pull/1111 for more details.
	KubernetesAPIApprovalPolicyConformant CustomResourceDefinitionConditionType = "KubernetesAPIApprovalPolicyConformant"
)

// CustomResourceDefinitionCondition contains details for the current condition of this pod.
type CustomResourceDefinitionCondition struct {
	// Type is the type of the condition. Types include Established, NamesAccepted and Terminating.
	Type CustomResourceDefinitionConditionType
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string
	// Human-readable message indicating details about last transition.
	// +optional
	Message string
}

// CustomResourceDefinitionStatus indicates the state of the CustomResourceDefinition
type CustomResourceDefinitionStatus struct {
	// Conditions indicate state for particular aspects of a CustomResourceDefinition
	// +listType=map
	// +listMapKey=type
	Conditions []CustomResourceDefinitionCondition

	// AcceptedNames are the names that are actually being used to serve discovery
	// They may be different than the names in spec.
	AcceptedNames CustomResourceDefinitionNames

	// StoredVersions are all versions of CustomResources that were ever persisted. Tracking these
	// versions allows a migration path for stored versions in etcd. The field is mutable
	// so the migration controller can first finish a migration to another version (i.e.
	// that no old objects are left in the storage), and then remove the rest of the
	// versions from this list.
	// None of the versions in this list can be removed from the spec.Versions field.
	StoredVersions []string
}

// CustomResourceCleanupFinalizer is the name of the finalizer which will delete instances of
// a CustomResourceDefinition
const CustomResourceCleanupFinalizer = "customresourcecleanup.apiextensions.k8s.io"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinition represents a resource that should be exposed on the API server.  Its name MUST be in the format
// <.spec.name>.<.spec.group>.
type AddonConfig struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec describes how the user wants the resources to appear
	Spec AddonConfigSpec
	// Status indicates the actual state of the CustomResourceDefinition
	Status CustomResourceDefinitionStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomResourceDefinitionList is a list of CustomResourceDefinition objects.
type AddonConfigList struct {
	metav1.TypeMeta
	metav1.ListMeta

	// Items individual CustomResourceDefinitions
	Items []AddonConfig
}

// CustomResourceValidation is a list of validation methods for CustomResources.
type CustomResourceValidation struct {
	// OpenAPIV3Schema is the OpenAPI v3 schema to be validated against.
	OpenAPIV3 *apiextensions.JSONSchemaProps
}

var _ resource.Object = &AddonConfig{}
var _ resourcestrategy.Validater = &AddonConfig{}

func (in *AddonConfig) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *AddonConfig) NamespaceScoped() bool {
	return true
}

func (in *AddonConfig) New() runtime.Object {
	return &AddonConfig{}
}

func (in *AddonConfig) NewList() runtime.Object {
	return &AddonConfigList{}
}

func (in *AddonConfig) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "addon.tanzu.vmware.com",
		Version:  "v1alpha1",
		Resource: "addonconfigs",
	}
}

func (in *AddonConfig) IsStorageVersion() bool {
	return true
}

func (in *AddonConfig) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &AddonConfigList{}

func (in *AddonConfigList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
