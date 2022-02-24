package deployment

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateDeployment(cluster string, node string, deployment *v1.Deployment) error {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	deployment.Spec.Template.Spec.NodeName = node
	deployment.ResourceVersion = ""

	namespace := deployment.ObjectMeta.Namespace
	CheckAndCreateNamespace(cluster, namespace)
	new_dep, err := cluster_client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return err
	} else {
		fmt.Printf("success to create %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
	}
	return nil
}
