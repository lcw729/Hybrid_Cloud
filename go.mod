module Hybrid_Cluster

go 1.16

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/aws/aws-sdk-go v1.40.29
	github.com/go-logr/logr v0.4.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/googleapis/gnostic v0.5.5
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.0.0-20211102192858-4dd72447c267 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5
	google.golang.org/api v0.54.0
	google.golang.org/genproto v0.0.0-20211102202547-e9cf271f7f2c // indirect
	google.golang.org/grpc v1.42.0
	// google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.4
	k8s.io/apiextensions-apiserver v0.22.4 // indirect
	k8s.io/apimachinery v0.22.4
	k8s.io/apiserver v0.22.4 // indirect
	k8s.io/autoscaler/vertical-pod-autoscaler v0.9.2
	k8s.io/client-go v0.22.4
	k8s.io/code-generator v0.22.4
	k8s.io/gengo v0.0.0-20201214224949-b6c5ce23f027
	k8s.io/klog/v2 v2.9.0
	k8s.io/kube-openapi v0.0.0-20211109043538-20434351676c
	k8s.io/sample-controller v0.16.8
	sigs.k8s.io/aws-iam-authenticator v0.5.3
	sigs.k8s.io/controller-runtime v0.9.6
	sigs.k8s.io/kubefed v0.8.1
)

replace (
	Hybrid_Cluster.com/policy-check v0.0.0 => ./policy-check
	admiralty.io/multicluster-controller => admiralty.io/multicluster-controller v0.1.0
	k8s.io/api => k8s.io/api v0.22.4
// k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.0-alpha.4
// k8s.io/apimachinery => k8s.io/apimachinery v0.23.0-alpha.4
// k8s.io/apiserver => k8s.io/apiserver v0.23.0-alpha.4
// k8s.io/cli-runtime => k8s.io/cli-runtime v0.23.0-alpha.4
// sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.0
)
