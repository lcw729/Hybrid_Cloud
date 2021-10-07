module Hybrid_Cluster

go 1.16

require (
	admiralty.io/multicluster-controller v0.6.0
	github.com/Azure/azure-sdk-for-go v58.0.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.40.29
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	google.golang.org/api v0.54.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/sample-controller v0.16.8
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/controller-runtime v0.9.6
	sigs.k8s.io/kubefed v0.8.1
)

replace (
	Hybrid_Cluster.com/policy-check v0.0.0 => ./policy-check
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.0
)
