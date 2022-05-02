package deployment

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetDeploymentName(pod *corev1.Pod) string {
	str := pod.GenerateName
	list := strings.SplitAfterN(str, "-", -1)
	length := len(list[len(list)-2]) + 1
	name := str[:len(str)-length]
	return name
}

func GetDeployment(cluster string, pod *corev1.Pod) (*v1.Deployment, error) {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/mnt/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	deploymentName := GetDeploymentName(pod)
	fmt.Println(deploymentName)
	d, err := cluster_client.AppsV1().Deployments(pod.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return d, err
	} else {
		fmt.Printf("success to get deployment %s in cluster %s [replicas : %d]\n", d.Name, cluster, *d.Spec.Replicas)
		return d, err
	}
}
