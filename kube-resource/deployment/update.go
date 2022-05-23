package deployment

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateDeployment(cluster string, deployment *v1.Deployment, replicas int32) error {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)

	if deployment.Status.Replicas == deployment.Status.ReadyReplicas && deployment.Status.Replicas == deployment.Status.AvailableReplicas {
		fmt.Println("deployment status is OK")
		deployment.Spec.Replicas = &replicas
		new_dep, err := cluster_client.AppsV1().Deployments("default").Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		} else {
			fmt.Printf("success to update %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
		}
	} else {
		fmt.Println("deployment status is not stable")
	}
	return nil
}
