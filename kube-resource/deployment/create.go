package deployment

import (
	"Hybrid_Cloud/hybridctl/util"
	cobrautil "Hybrid_Cloud/hybridctl/util"
	resourcev1alpha1apis "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"context"
	"fmt"

	ns "Hybrid_Cloud/kube-resource/namespace"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/apps/v1"
)

func CreateDeployment(cluster string, node string, deployment *v1.Deployment) error {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	deployment.Spec.Template.Spec.NodeName = node
	deployment.ResourceVersion = ""

	namespace := deployment.ObjectMeta.Namespace
	ns.CreateNamespace(cluster, namespace)
	new_dep, err := cluster_client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return err
	} else {
		fmt.Printf("success to create %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
	}
	return nil
}

func DeployDeploymentFromHCPDeployment(hcp_resource *resourcev1alpha1apis.HCPDeployment) bool {
	targets := hcp_resource.Spec.SchedulingResult.Targets
	metadata := hcp_resource.Spec.RealDeploymentMetadata
	if metadata.Namespace == "" {
		metadata.Namespace = "default"
	}
	spec := hcp_resource.Spec.RealDeploymentSpec

	// HCPDeployment SchedulingResult에 따라 Deployment 배포
	for _, target := range targets {
		// cluster clientset 생성

		config, err := util.BuildConfigFromFlags(target.Cluster, "/root/.kube/config")
		if err != nil {
			fmt.Println(err)
			return false
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			return false
		}

		// spec 값 재설정하기
		spec.Replicas = target.Replicas

		// 배포할 Deployment resource 정의
		kube_resource := appsv1.Deployment{
			ObjectMeta: metadata,
			Spec:       spec,
		}

		// Deployment 배포
		r, err := clientset.AppsV1().Deployments(metadata.Namespace).Create(context.TODO(), &kube_resource, metav1.CreateOptions{})

		if err != nil {
			fmt.Println(err)
			return false
		} else {
			fmt.Printf("succeed to deploy deployment %s in %s\n", r.ObjectMeta.Name, target.Cluster)
		}
	}
	return true
}
