package namespace

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckAndCreateNamespace(cluster string, namespace string) (*corev1.Namespace, error) {
	config, err := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
		ns, err := client.CoreV1().Namespaces().Create(context.TODO(), &Namespace, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			fmt.Printf("success to create namespace %s in %s", namespace, cluster)
			return ns, nil
		}
	}
	return nil, err
}
