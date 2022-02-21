package resource

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckAndCreateNamespace(cluster string, namespace string) {
	config, err := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	if !FindNamespaceList(cluster, namespace) {
		Namespace := corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = client.CoreV1().Namespaces().Create(context.TODO(), &Namespace, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("success to create namespace %s in %s", namespace, cluster)
		}
	}
}

func CreateDeployment(cluster string, node string, deployment *v1.Deployment) error {
	fmt.Println("111111")
	config, err := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	fmt.Println("22222")
	if err != nil {
		fmt.Println("this error")
		return err
	}
	cluster_client := kubernetes.NewForConfigOrDie(config)
	deployment.Spec.Template.Spec.NodeName = node
	deployment.ResourceVersion = ""

	namespace := deployment.ObjectMeta.Namespace
	CheckAndCreateNamespace(cluster, namespace)
	new_dep, err := cluster_client.AppsV1().Deployments(deployment.ObjectMeta.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("success to create %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
	}
	return nil
}
