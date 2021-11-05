package v1alpha1

import (
	policyv1alpha1 "Hybrid_Cluster/apis/policy/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HCPDeploymentSpec defines the desired state of HCPDeployment
// +k8s:openapi-gen=true
type HCPDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// changes
	Template HCPDeploymentTemplate `json:"template" protobuf:"bytes,3,opt,name=template"`

	// Added
	Replicas int32               `json:"replicas" protobuf:"varint,1,opt,name=replicas"`
	Labels   map[string]string   `json:"labels,omitempty" protobuf:"bytes,11,opt,name=labels"`
	Affinity map[string][]string `json:"affinity,omitempty" protobuf:"bytes,3,opt,name=affinity"`
	Policy   map[string]string   `json:"policy,omitempty" protobuf:"bytes,3,opt,name=policy"`
	//Placement
}

type HCPDeploymentTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// changed
	Spec   HCPDeploymentTemplateSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status appsv1.DeploymentStatus   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type HCPDeploymentTemplateSpec struct {
	Replicas *int32                `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
	// changed
	Template                HCPPodTemplateSpec        `json:"template" protobuf:"bytes,3,opt,name=template"`
	Strategy                appsv1.DeploymentStrategy `json:"strategy,omitempty" patchStrategy:"retainKeys" protobuf:"bytes,4,opt,name=strategy"`
	MinReadySeconds         int32                     `json:"minReadySeconds,omitempty" protobuf:"varint,5,opt,name=minReadySeconds"`
	RevisionHistoryLimit    *int32                    `json:"revisionHistoryLimit,omitempty" protobuf:"varint,6,opt,name=revisionHistoryLimit"`
	Paused                  bool                      `json:"paused,omitempty" protobuf:"varint,7,opt,name=paused"`
	ProgressDeadlineSeconds *int32                    `json:"progressDeadlineSeconds,omitempty" protobuf:"varint,9,opt,name=progressDeadlineSeconds"`
}

type HCPPodTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// changed
	Spec HCPPodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type HCPPodSpec struct {
	Volumes []corev1.Volume `json:"volumes,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,1,rep,name=volumes"`
	// changes
	InitContainers []HCPContainer `json:"initContainers,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,20,rep,name=initContainers"`
	// changes
	Containers                    []HCPContainer                `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`
	RestartPolicy                 corev1.RestartPolicy          `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy,casttype=RestartPolicy"`
	TerminationGracePeriodSeconds *int64                        `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	ActiveDeadlineSeconds         *int64                        `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,5,opt,name=activeDeadlineSeconds"`
	DNSPolicy                     corev1.DNSPolicy              `json:"dnsPolicy,omitempty" protobuf:"bytes,6,opt,name=dnsPolicy,casttype=DNSPolicy"`
	NodeSelector                  map[string]string             `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	ServiceAccountName            string                        `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`
	DeprecatedServiceAccount      string                        `json:"serviceAccount,omitempty" protobuf:"bytes,9,opt,name=serviceAccount"`
	AutomountServiceAccountToken  *bool                         `json:"automountServiceAccountToken,omitempty" protobuf:"varint,21,opt,name=automountServiceAccountToken"`
	NodeName                      string                        `json:"nodeName,omitempty" protobuf:"bytes,10,opt,name=nodeName"`
	HostNetwork                   bool                          `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	HostPID                       bool                          `json:"hostPID,omitempty" protobuf:"varint,12,opt,name=hostPID"`
	HostIPC                       bool                          `json:"hostIPC,omitempty" protobuf:"varint,13,opt,name=hostIPC"`
	ShareProcessNamespace         *bool                         `json:"shareProcessNamespace,omitempty" protobuf:"varint,27,opt,name=shareProcessNamespace"`
	SecurityContext               *corev1.PodSecurityContext    `json:"securityContext,omitempty" protobuf:"bytes,14,opt,name=securityContext"`
	ImagePullSecrets              []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`
	Hostname                      string                        `json:"hostname,omitempty" protobuf:"bytes,16,opt,name=hostname"`
	Subdomain                     string                        `json:"subdomain,omitempty" protobuf:"bytes,17,opt,name=subdomain"`
	Affinity                      *corev1.Affinity              `json:"affinity,omitempty" protobuf:"bytes,18,opt,name=affinity"`
	SchedulerName                 string                        `json:"schedulerName,omitempty" protobuf:"bytes,19,opt,name=schedulerName"`
	Tolerations                   []corev1.Toleration           `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	HostAliases                   []corev1.HostAlias            `json:"hostAliases,omitempty" patchStrategy:"merge" patchMergeKey:"ip" protobuf:"bytes,23,rep,name=hostAliases"`
	PriorityClassName             string                        `json:"priorityClassName,omitempty" protobuf:"bytes,24,opt,name=priorityClassName"`
	Priority                      *int32                        `json:"priority,omitempty" protobuf:"bytes,25,opt,name=priority"`
	DNSConfig                     *corev1.PodDNSConfig          `json:"dnsConfig,omitempty" protobuf:"bytes,26,opt,name=dnsConfig"`
	ReadinessGates                []corev1.PodReadinessGate     `json:"readinessGates,omitempty" protobuf:"bytes,28,opt,name=readinessGates"`
	RuntimeClassName              *string                       `json:"runtimeClassName,omitempty" protobuf:"bytes,29,opt,name=runtimeClassName"`
	EnableServiceLinks            *bool                         `json:"enableServiceLinks,omitempty" protobuf:"varint,30,opt,name=enableServiceLinks"`
}

type HCPContainer struct {
	Name       string                 `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image      string                 `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Command    []string               `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args       []string               `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir string                 `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	Ports      []corev1.ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
	EnvFrom    []corev1.EnvFromSource `json:"envFrom,omitempty" protobuf:"bytes,19,rep,name=envFrom"`
	Env        []corev1.EnvVar        `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
	// changes
	Resources                HCPResourceRequirements         `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	VolumeMounts             []corev1.VolumeMount            `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	VolumeDevices            []corev1.VolumeDevice           `json:"volumeDevices,omitempty" patchStrategy:"merge" patchMergeKey:"devicePath" protobuf:"bytes,21,rep,name=volumeDevices"`
	LivenessProbe            *corev1.Probe                   `json:"livenessProbe,omitempty" protobuf:"bytes,10,opt,name=livenessProbe"`
	ReadinessProbe           *corev1.Probe                   `json:"readinessProbe,omitempty" protobuf:"bytes,11,opt,name=readinessProbe"`
	Lifecycle                *corev1.Lifecycle               `json:"lifecycle,omitempty" protobuf:"bytes,12,opt,name=lifecycle"`
	TerminationMessagePath   string                          `json:"terminationMessagePath,omitempty" protobuf:"bytes,13,opt,name=terminationMessagePath"`
	TerminationMessagePolicy corev1.TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty" protobuf:"bytes,20,opt,name=terminationMessagePolicy,casttype=TerminationMessagePolicy"`
	ImagePullPolicy          corev1.PullPolicy               `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	SecurityContext          *corev1.SecurityContext         `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`
	Stdin                    bool                            `json:"stdin,omitempty" protobuf:"varint,16,opt,name=stdin"`
	StdinOnce                bool                            `json:"stdinOnce,omitempty" protobuf:"varint,17,opt,name=stdinOnce"`
	TTY                      bool                            `json:"tty,omitempty" protobuf:"varint,18,opt,name=tty"`
}

type HCPResourceRequirements struct {
	Limits   corev1.ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	Requests corev1.ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
	// Added
	Needs isNeedResourceList `json:"needs,omitempty" protobuf:"bytes,2,rep,name=needs,casttype=isNeedResourceList,castkey=ResourceName"`
}

type isNeedResourceList map[corev1.ResourceName]bool

// HCPDeploymentStatus defines the observed state of HCPDeployment
// +k8s:openapi-gen=true
type HCPDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas                  int32             `json:"replicas"`
	ClusterMaps               map[string]int32  `json:"clusters"`
	LastSpec                  HCPDeploymentSpec `json:"lastSpec"`
	SchedulingNeed            bool              `json:"schedulingNeed"`
	SchedulingComplete        bool              `json:"schedulingComplete"`
	CreateSyncRequestComplete bool              `json:"createSyncRequestComplete"`
	SyncRequestName           string            `json:"syncRequestName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPDeployment is the Schema for the HCPdeployments API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   HCPDeploymentSpec   `json:"spec,omitempty"`
	Status HCPDeploymentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPDeploymentList contains a list of HCPDeployment
type HCPDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPDeployment `json:"items"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HCPIngressSpec defines the desired state of HCPIngress
// +k8s:openapi-gen=true
type HCPIngressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template extv1b1.Ingress `json:"template" protobuf:"bytes,3,opt,name=template"`
	//Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	//Placement

}

// HCPIngressStatus defines the observed state of HCPIngress
// +k8s:openapi-gen=true
type HCPIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// Replicas int32 `json:"replicas"`
	ClusterMaps map[string]int32 `json:"clusters"`
	ChangeNeed  bool             `json:"changeNeed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPIngress is the Schema for the HCPingresss API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPIngressSpec   `json:"spec,omitempty"`
	Status HCPIngressStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPIngressList contains a list of HCPIngress
type HCPIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPIngress `json:"items"`
}

// HCPServiceSpec defines the desired state of HCPService
// +k8s:openapi-gen=true
type HCPServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	LabelSelector map[string]string `json:"labelselector" protobuf:"bytes,1,opt,name=labelselector"`
	Template      corev1.Service    `json:"template" protobuf:"bytes,2,opt,name=template"`
	//Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	//Placement

}

// HCPServiceStatus defines the observed state of HCPService
// +k8s:openapi-gen=true
type HCPServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	//Replicas int32 `json:"replicas"`
	ClusterMaps map[string]int32 `json:"clusters"`
	LastSpec    HCPServiceSpec   `json:"lastSpec"`
	ChangeNeed  bool             `json:"changeNeed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPService is the Schema for the HCPservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPServiceSpec   `json:"spec,omitempty"`
	Status HCPServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPServiceList contains a list of HCPService
type HCPServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPService `json:"items"`
}

type HCPHybridAutoScalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	MainController string         `json:"mainController" protobuf:"bytes,1,opt,name=maincontroller"`
	ScalingOptions ScalingOptions `json:"scalingOptions,omitempty" protobuf:"bytes,2,opt,name=scalingoptions"`
}

type ScalingOptions struct {
	CpaTemplate CpaTemplate                        `json:"cpaTemplate,omitempty" protobuf:"bytes,1,opt,name=cpatemplate"`
	HpaTemplate hpav2beta1.HorizontalPodAutoscaler `json:"hpaTemplate,omitempty" protobuf:"bytes,2,opt,name=hpatemplate"`
	VpaTemplate string                             `json:"vpaTemplate,omitempty" protobuf:"bytes,3,opt,name=vpatemplate"`
	//VpaTemplate vpav1beta2.VerticalPodAutoscaler `json:"vpaTemplate" protobuf:"bytes,3,opt,name=vpaTemplate"`
}

type CpaTemplate struct {
	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaletargetref"`
	MinReplicas    int32          `json:"minReplicas" protobuf:"varint,2,opt,name=minreplicas"`
	MaxReplicas    int32          `json:"maxReplicas" protobuf:"varint,3,opt,name=maxreplicas"`
}

type ScaleTargetRef struct {
	// Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
}

// HCPHybridAutoScalerStatus defines the observed state of HCPHybridAutoScaler
// +k8s:openapi-gen=true
type HCPHybridAutoScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	//Nodes []string `json:"nodes"`
	LastSpec         HCPHybridAutoScalerSpec      `json:"lastSpec"`
	Policies         []policyv1alpha1.HCPPolicies `json:"policies"`
	RebalancingCount map[string]int32             `json:"rebalancingCount"`
	SyncRequestName  string                       `json:"syncRequestName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPHybridAutoScaler is the Schema for the HCPhybridautoscalers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPHybridAutoScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPHybridAutoScalerSpec   `json:"spec,omitempty"`
	Status HCPHybridAutoScalerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPHybridAutoScalerList contains a list of HCPHybridAutoScaler
type HCPHybridAutoScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPHybridAutoScaler `json:"items"`
}

// HCPConfigMapSpec defines the desired state of HCPConfigMap
// +k8s:openapi-gen=true
type HCPConfigMapSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template corev1.ConfigMap `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// HCPConfigMapStatus defines the observed state of HCPConfigMap
// +k8s:openapi-gen=true
type HCPConfigMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps     map[string]int32 `json:"clusters"`
	SyncRequestName string           `json:"syncRequestName"`
}

// HCPConfigMap is the Schema for the HCPconfigmaps API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPConfigMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPConfigMapSpec   `json:"spec,omitempty"`
	Status HCPConfigMapStatus `json:"status,omitempty"`
}

// HCPConfigMapList contains a list of HCPConfigMap
type HCPConfigMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPConfigMap `json:"items"`
}

type HCPSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template corev1.Secret `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// HCPSecretStatus defines the observed state of HCPSecret
// +k8s:openapi-gen=true
type HCPSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps map[string]int32 `json:"clusters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPSecret is the Schema for the HCPsecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPSecretSpec   `json:"spec,omitempty"`
	Status HCPSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPSecretList contains a list of HCPSecret
type HCPSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPSecret `json:"items"`
}

type HCPJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template batchv1.Job `json:"template" protobuf:"bytes,3,opt,name=template"`
}

// HCPJobStatus defines the observed state of HCPJob
// +k8s:openapi-gen=true
type HCPJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps map[string]int32 `json:"clusters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPJob is the Schema for the HCPjobs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPJobSpec   `json:"spec,omitempty"`
	Status HCPJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPJobList contains a list of HCPJob
type HCPJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPJob `json:"items"`
}

type HCPNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template appsv1.Deployment `json:"template" protobuf:"bytes,3,opt,name=template"`
	//Placement
}

// HCPNamespaceStatus defines the observed state of HCPNamespace
// +k8s:openapi-gen=true
type HCPNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ClusterMaps map[string]int32 `json:"clusters"`
	ChangeNeed  bool             `json:"changeneed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPNamespace is the Schema for the HCPnamespaces API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPNamespaceSpec   `json:"spec,omitempty"`
	Status HCPNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPNamespaceList contains a list of HCPNamespace
type HCPNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPNamespace `json:"items"`
}

//func init() {
//	SchemeBuilder.Register(&HCPDeployment{}, &HCPDeploymentList{})
//	SchemeBuilder.Register(&HCPIngress{}, &HCPIngressList{})
//	SchemeBuilder.Register(&HCPService{}, &HCPServiceList{})
//	SchemeBuilder.Register(&HCPHybridAutoScaler{}, &HCPHybridAutoScalerList{})
//	SchemeBuilder.Register(&HCPPolicyEngine{}, &HCPPolicyEngineList{})
//
//}
