package namespace

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FindNamespaceList(cluster string, namespace string) bool {

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
			return true
		}
	}
	return false
}
