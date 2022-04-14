package namespace

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNamespace(cluster string, namespace string) *corev1.Namespace {

	config, err := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	namespaceList, _ := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	for i := range namespaceList.Items {

		if namespaceList.Items[i].Name == namespace {
			return &namespaceList.Items[i]
		}
	}
	return nil
}
