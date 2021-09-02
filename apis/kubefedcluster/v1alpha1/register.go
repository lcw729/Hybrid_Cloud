package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "KubeFedCluster"
	// GroupVersion is the version.
	GroupVersion string = "v1alpha1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "KubeFedCluster"
	// Singular is used as an alias on kubectl for display.
	Singular string = "KubeFedCluster"
	// CRDName is the CRD name for Jinghzhu.
	CRDName string = Plural + "." + "core.kubefed.io"
	// ShortName is the short alias for the CRD.
	ShortName string = "kfc"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{
		Group:   fedv1b1.SchemeGroupVersion.Group,
		Version: fedv1b1.SchemeGroupVersion.Version,
	}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&fedv1b1.KubeFedCluster{},
		&fedv1b1.KubeFedClusterList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
