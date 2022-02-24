package pod

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
)

func GetPod(cluster string, pod string, pod_namespace string) (*corev1.Pod, error) {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	p, err := cluster_client.CoreV1().Pods(pod_namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if err != nil {
		return p, err
	} else {
		fmt.Printf("success to get pod %s [cluster %s, node %s]\n", p.Name, cluster, p.Spec.NodeName)
		return p, err
	}
}
